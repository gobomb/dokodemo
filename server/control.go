package server

import (
	"doko/conn"
	"doko/msg"
	"doko/util"
	"fmt"
	"github.com/qiniu/log"
	"io"
	"runtime/debug"
	"strings"
	"time"
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

	randStr := util.GenerateRandomString()
	c.id = authMsg.ClientId
	if c.id == "" {
		c.id = randStr
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
	go c.stopper()

}

func (c *Control) registerTunnel(rawTunnelReq *msg.ReqTunnel) {
	//for _,
	tunnelReq := *rawTunnelReq

	log.Printf("Registering new tunnel")
	t := NewTunnel(&tunnelReq, c)
	c.tunnels = append(c.tunnels, t)
	c.out <- &msg.NewTunnel{
		Url:      t.url,
		Protocol: "tcp",
		ReqId:    rawTunnelReq.ReqId,
	}
	log.Info(rawTunnelReq.ReqId)
	rawTunnelReq.Hostname = strings.Replace(t.url, "tcp"+"://", "", 1)
}

func (c *Control) manager() {
	//defer func() {
	//	if err := recover(); err != nil {
	//		log.Printf("Control::manager failed with error %v: %s", err, debug.Stack())
	//	}
	//}()
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
				log.Debugf("msg.ReqTunnel:%v", m)
				c.registerTunnel(m)
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
			} else {
				log.Errorf("[msg.ReadMsg] %v", err)
				//(err)
			}
		} else {
			c.in <- msg
		}
	}

}

func (c *Control) GetProxy() (conn.Conn, error) {
	var (
		proxyConn conn.Conn
		err       error
	)
	var ok bool
	//var err error
	select {
	case proxyConn, ok = <-c.proxies:
		if !ok {
			err = fmt.Errorf("no proxy connections available, control is closing")
			return nil, err
		}
	default:
		log.Printf("No proxy in pool, requesting proxy from control . . .")
		c.out <- &msg.ReqProxy{}
		select {
		case proxyConn, ok = <-c.proxies:
			if !ok {
				err = fmt.Errorf("no proxy connections available, control is closing")
				return nil, err
			}
		case <-time.After(pingTimeoutInterval):
			err = fmt.Errorf("timeout tring to get proxy connection")
			return nil, err
		}
	}
	return proxyConn, nil

}

func (c *Control) RegisterProxy(conn conn.Conn) {
	conn.SetDeadline(time.Now().Add(proxyStaleDuration))
	select {
	case c.proxies <- conn:
		log.Printf("Registered")

	default:
		log.Printf("Proxies buffer is full, discarding.")
		conn.Close()
	}
}

func (c *Control) stopper() {
	c.shutdown.WaitBegin()

	controlRegistry.Del(c.id)

	close(c.in)

	c.managerShutdown.WaitComplete()

	close(c.out)
	c.conn.Close()
	for _, t := range c.tunnels {
		t.Shutdown()
	}

	close(c.proxies)

	for p := range c.proxies {
		p.Close()
	}

	c.shutdown.Complete()
	log.Printf("Shutdown complete")
}
