package main

import (
	"github.com/qiniu/log"
	"sync"
	"time"
	"net"
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

type Conn interface {
	net.Conn
	Id() string
	SetType(string)
	CloseRead() error
}

type Message interface{}

type Shutdown struct {
	sync.Mutex
	inProgress bool
	begin      chan int // closed when the shutdown begins
	complete   chan int // closed when the shutdown completes
}

type Control struct {
	// actual connection
	conn Conn

	// put a message in this channel to send it over
	// conn to the client
	out chan Message

	// read from this channel to get the next message sent
	// to us over conn by the client
	in chan Message

	// the last time we received a ping from the client - for heartbeats
	lastPing time.Time

	// all of the tunnels this control connection handles
	tunnels []*Tunnel

	// proxy connections
	proxies chan Conn

	// identifier
	id string

	// synchronizer for controlled shutdown of writer()
	writerShutdown *Shutdown

	// synchronizer for controlled shutdown of reader()
	readerShutdown *Shutdown

	// synchronizer for controlled shutdown of manager()
	managerShutdown *Shutdown

	// synchronizer for controller shutdown of entire Control
	shutdown *Shutdown
}

type Tunnel struct {
	// 隧道建立请求
	req *ReqTunnel

	// time when the tunnel was opened
	start time.Time

	// 公网 url
	url string

	// tcp 监听
	listener *net.TCPListener

	// 控制连接
	ctl *Control

	// closing
	closing int32
}

type TunnelRegistry struct {
	tunnels map[string]*Tunnel
	//affinity *cache.LRUCache
	sync.RWMutex
}
type ControlRegistry struct {
	controls map[string]*Control
	sync.RWMutex
}
type loggedConn struct {
	tcp *net.TCPConn
	net.Conn
	id  int32
	typ string
}

type Listener struct {
	net.Addr
	Conns chan *loggedConn
}

var (
	tunnelRegistry  *TunnelRegistry
	controlRegistry *ControlRegistry

	// XXX: kill these global variables - they're only used in tunnel.go for constructing forwarding URLs
	opts      *Options
	listeners map[string]*Listener
)

func main() {
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
	listeners = make(map[string]*Listener)
	tunnelListener(opts.tunnelAddr)
}

func tunnelListener(addr string){

}
