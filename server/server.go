package server

import (
	"github.com/qiniu/log"
	"time"
	"runtime/debug"
	"doko/conn"
	"doko/msg"
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

//type TunnelRegistry struct {
//	tunnels map[string]*Tunnel
//	//affinity *cache.LRUCache
//	sync.RWMutex
//}
//type ControlRegistry struct {
//	controls map[string]*Control
//	sync.RWMutex
//}

var (
	tunnelRegistry  *TunnelRegistry
	controlRegistry *ControlRegistry

	// XXX: kill these global variables - they're only used in tunnel.go for constructing forwarding URLs
	opts      *Options
	listeners map[string]*conn.Listener
)

func Main() {
	log.Print("start!")
	opts := &Options{
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

			// don't timeout after the initial read, tunnel heartbeating will kill
			// dead connections
			tunnelConn.SetReadDeadline(time.Time{})

			switch m := rawMsg.(type) {
			case *msg.Auth:
				NewControl(tunnelConn, m)

				//case *msg.RegProxy:
				//	NewProxy(tunnelConn, m)

			default:
				tunnelConn.Close()
			}
			//rawMsg,err:=msg.ReadMsg(tunnelConn)
			//if err!=nil{
			//	log.Printf("[msg.ReadMsg] %v",err)
			//}
			//log.Printf("[rawMsg] %v",rawMsg)

		}(c)
	}
}
