package server

import (
	"doko/conn"
	"doko/msg"
	"github.com/qiniu/log"
	"runtime/debug"
	"time"
)

const (
	connReadTimeout = 10 * time.Second
)

type Options struct {
	tunnelAddr string
	domain     string
}

type ReqTunnel struct {
	ReqId string

	// tcp only
	RemotePort uint16
}

type Info struct{
	TunReg  *TunnelRegistry
	CtlReg *ControlRegistry
	Opts      *Options
	Listeners map[string]*conn.Listener
}
var (
	tunnelRegistry  *TunnelRegistry
	controlRegistry *ControlRegistry

	// XXX: kill these global variables - they're only used in tunnel.go for constructing forwarding URLs
	opts      *Options
	listeners map[string]*conn.Listener
)

func GetInfo() Info{
	return Info{
		tunnelRegistry,
		controlRegistry,
		opts,
		listeners,
	}
}

func Main() {
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
	tunnelListener(opts.tunnelAddr)
}

func tunnelListener(addr string) {

	listener := conn.Listen(addr)

	log.Printf("Listening for control and proxy connections on %s", listener.Addr.String())
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
