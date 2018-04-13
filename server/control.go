package server

import (
	"github.com/qiniu/log"
	"time"
	"doko/conn"
	"doko/msg"
	"doko/util"
	"runtime/debug"
	"io"
)

const (
	pingTimeoutInterval = 30 * time.Second
	connReapInterval    = 10 * time.Second
	controlWriteTimeout = 10 * time.Second
	proxyStaleDuration  = 60 * time.Second
	proxyMaxPoolSize    = 10
)

type Control struct {
	auth *msg.Auth
	// actual connection
	conn conn.Conn

	// put a message in this channel to send it over
	// conn to the client
	out chan msg.Message

	// read from this channel to get the next message sent
	// to us over conn by the client
	in chan msg.Message

	// the last time we received a ping from the client - for heartbeats
	lastPing time.Time

	// all of the tunnels this control connection handles
	tunnels []*Tunnel

	// proxy connections
	proxies chan conn.Conn

	// identifier
	id string

	// synchronizer for controlled shutdown of writer()
	writerShutdown *util.Shutdown

	// synchronizer for controlled shutdown of reader()
	readerShutdown *util.Shutdown

	// synchronizer for controlled shutdown of manager()
	managerShutdown *util.Shutdown

	// synchronizer for controller shutdown of entire Control
	shutdown *util.Shutdown
}

func NewControl(ctlConn conn.Conn, authMsg *msg.Auth) {
	//var err error

	// 创建控制器实例
	c := &Control{
		auth:            authMsg,
		conn:            ctlConn,
		out:             make(chan msg.Message),
		in:              make(chan msg.Message),
		proxies:         make(chan conn.Conn, 10),
		lastPing:        time.Now(),
		writerShutdown:  util.NewShutdown(),
		readerShutdown:  util.NewShutdown(),
		managerShutdown: util.NewShutdown(),
		shutdown:        util.NewShutdown(),
	}

	// 错误返回匿名函数
	//failAuth := func(e error) {
	//	_ = msg.WriteMsg(ctlConn, &msg.AuthResp{Error: e.Error{}})
	//	ctlConn.Close()
	//}

	// 登记客户端 id
	c.id = authMsg.ClientId
	if c.id == ""{
		c.id = "1111"
	}

	// 注册控制器
	if replaced := controlRegistry.Add(c.id, c); replaced != nil {
		replaced.shutdown.WaitComplete()
	}

	go c.writer()

	c.out <- &msg.AuthResp{
		ClientId: c.id,
	}

	c.out <- &msg.ReqProxy{}

	go c.manager()
	go c.reader()
	//go c.stopper()

}

func (c *Control) manager() {
	defer func() {
		if err := recover(); err != nil {
			log.Printf("Control::manager failed with error %v: %s", err, debug.Stack())
		}
	}()
	defer c.shutdown.Begin()

	defer c.managerShutdown.Complete()

	reap := time.NewTicker(connReapInterval)
	defer reap.Stop()

	for {
		select {
		case <-reap.C:
			if time.Since(c.lastPing) > pingTimeoutInterval {
				log.Printf("Lost heartbeat")
				c.shutdown.Begin()
			}

		case mRaw, ok := <-c.in:
			if !ok {
				return
			}
			switch m := mRaw.(type) {
			case *msg.ReqTunnel:
				log.Println(m)
				//c.registerTunnel(m)
			case *msg.Ping:
				c.lastPing = time.Now()
				c.out <- &msg.Pong{}
			}
		}
	}

}

func (c *Control) Replaced(replacement *Control) {
	log.Printf("Replaced by control: %s", replacement.conn.Id())
	c.id = ""
	c.shutdown.Begin()
}

func (c *Control) writer() {
	defer func() {
		if err := recover(); err != nil {
			log.Printf("Control::writer failed with error %v: %s", err, debug.Stack())
		}
	}()

	defer c.shutdown.Begin()

	defer c.writerShutdown.Complete()

	for m := range c.out {
		c.conn.SetWriteDeadline(time.Now().Add(controlWriteTimeout))
		if err := msg.WriteMsg(c.conn, m); err != nil {
			log.Printf("[msg.writeMsg] %v", err)
			panic(err)
		}
	}
}

func (c *Control) reader() {
	defer func() {
		if err := recover(); err != nil {
			log.Printf("Control::reader failed with error %v: %s", err, debug.Stack())
		}
	}()
	defer c.shutdown.Begin()

	defer c.readerShutdown.Complete()

	for {
		if msg, err := msg.ReadMsg(c.conn); err != nil {
			if err == io.EOF {
				log.Println("EOF")
				return
			}else{
				log.Printf("[msg.ReadMsg] %v",err)
				panic(err)
			}
		} else{
			c.in<-msg
		}
	}

}
