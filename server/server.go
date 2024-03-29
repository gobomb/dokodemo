package server

import (
	"doko/conn"
	"doko/msg"
	"github.com/qiniu/log"
	"runtime/debug"
	"time"
	"sync"
)

const (
	connReadTimeout = 10 * time.Second
)
type Status struct{
	StatusOn bool
	sync.Mutex
	}
var S = Status{
	StatusOn:false,
}
//var

var StopChan chan interface{}

type Options struct {
	tunnelAddr string
	domain     string
}

type ReqTunnel struct {
	ReqId string

	// tcp only
	RemotePort uint16
}

var (
	tunnelRegistry  *TunnelRegistry
	controlRegistry *ControlRegistry

	// XXX: kill these global variables - they're only used in tunnel.go for constructing forwarding URLs
	opts      *Options
	listeners map[string]*conn.Listener
)



func Main(stopChan chan interface{}) {
	StopChan = stopChan
	log.Print("start!")
	opts = &Options{
		tunnelAddr: ":4443",
		domain:     "127.0.0.1",
	}
	tunnelRegistry = &TunnelRegistry{
		tunnels: make(map[string]*Tunnel),
	}
	controlRegistry = &ControlRegistry{
		controls: make(map[string]*Control),
	}
	listeners = make(map[string]*conn.Listener)
	S.Lock()
	S.StatusOn=true
	S.Unlock()
	tunnelListener(opts.tunnelAddr)
}

func tunnelListener(addr string) {

	listener, err := conn.Listen(addr)
	if err != nil {
		log.Errorf("listen err: %v", err)
		return
	}
	log.Printf("Listening for control and proxy connections on %s", listener.Addr.String())

	// 等待来自 http 网关的停止消息
	go func() {
		<-StopChan
		// 关闭与各个客户端的代理连接
		for k, v := range controlRegistry.controls {
			log.Printf("Closing control %v", k)
			v.shutdown.Begin()
			v.shutdown.WaitComplete()
		}
		//for k,v:=range tunnelRegistry.tunnels{
		//	log.Infof("Closing tunnel listener: %v",k)
		//	v.listener.Close()
		//}
		// 将服务状态更改为false
		S.Lock()
		S.StatusOn=false
		S.Unlock()
		//StatusOn = false
		// 关闭监听
		//listener.Shutdown.Begin()
	}()

	for c := range listener.Conns {
		go func(tunnelConn conn.Conn) {
			defer func() {
				if r := recover(); r != nil {
					log.Printf("tunnelListener failed with error %v: %s", r, debug.Stack())
				}
			}()
			tunnelConn.SetReadDeadline(time.Now().Add(connReadTimeout))
			var rawMsg msg.Message
			rawMsg, err := msg.ReadMsg(tunnelConn)
			if err != nil {
				log.Printf("Failed to read message: %v", err)
				tunnelConn.Close()
				return
			}

			// don't timeout after the initial read, tunnel heartbeating will kill dead connections
			tunnelConn.SetReadDeadline(time.Time{})

			switch m := rawMsg.(type) {
			case *msg.Auth:
				NewControl(tunnelConn, m)

			case *msg.RegProxy:
				NewProxy(tunnelConn, m)

			default:
				tunnelConn.Close()
			}

		}(c)
	}
}

func NewProxy(pxyConn conn.Conn, regPxy *msg.RegProxy) {
	// fail gracefully if the proxy connection fails to register
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Failed with error: %v", r)
			pxyConn.Close()
		}
	}()

	// set logging prefix
	pxyConn.SetType("pxy")

	// look up the control connection for this proxy
	log.Printf("Registering new proxy for %s", regPxy.ClientId)
	ctl := controlRegistry.Get(regPxy.ClientId)

	if ctl == nil {
		panic("No client found for identifier: " + regPxy.ClientId)
	}

	ctl.RegisterProxy(pxyConn)
}
