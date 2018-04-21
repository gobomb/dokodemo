package server

import (
	"doko/conn"
	"doko/msg"
	"fmt"
	"github.com/qiniu/log"
	"net"
	"sync/atomic"
	"time"
)

type Tunnel struct {
	// 隧道建立请求
	req *msg.ReqTunnel

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

func (t *Tunnel) Shutdown() {
	log.Printf("Shutting down")
	atomic.StoreInt32(&t.closing, 1)
	if t.listener != nil {
		t.listener.Close()
	}
	tunnelRegistry.Del(t.url)
}

func NewTunnel(m *msg.ReqTunnel, ctl *Control) (t *Tunnel) {
	var err error
	t = &Tunnel{
		req:   m,
		start: time.Now(),
		ctl:   ctl,
	}

	proto := t.req.Protocol
	switch proto {
	case "tcp":
		bindTcp := func(port int) {
			t.listener, err = net.ListenTCP("tcp", &net.TCPAddr{IP: net.ParseIP("0.0.0.0"), Port: port})
			if err != nil {
				err = fmt.Errorf("Error binding TCP listener: %v\n", err)
				log.Println(err)
				//return err
			}
			addr := t.listener.Addr().(*net.TCPAddr)
			t.url = fmt.Sprintf("tcp://%s:%d", opts.domain, addr.Port)

			if err = tunnelRegistry.Register(t.url, t); err != nil {
				t.listener.Close()
				log.Printf("%v", err)
				err = fmt.Errorf("TCP listener bound, but failed to register %s", t.url)
				log.Println(err)
			}
			go t.listenTcp(t.listener)
		}

		if t.req.RemotePort != 0 {
			bindTcp(int(t.req.RemotePort))
			return
		}

		bindTcp(0)
		return
	default:
		err = fmt.Errorf("protocol %s is not supported", proto)
		return
	}
	log.Printf("Registerd new tunnel on: %s", t.ctl.conn.Id())
	return t
}

func (t *Tunnel) listenTcp(listener *net.TCPListener) {
	for {
		//var tcpConn conn.Conn
		tcpConn, err := listener.AcceptTCP()
		if err != nil {
			if atomic.LoadInt32(&t.closing) == 1 {
				return
			}
			log.Printf("Failed to accept new TCP conn: %v", err)
			continue
		}
		wrappedConn := conn.Wrap(tcpConn, "pub")
		log.Printf("New connection from %v", tcpConn.RemoteAddr())
		go t.HandlePublicConnection(wrappedConn)
	}

}

func (t *Tunnel) HandlePublicConnection(publicConn conn.Conn) {
	startTime := time.Now()
	var proxyConn conn.Conn
	var err error

	for i := 0; i < (2 * proxyMaxPoolSize); i++ {

		if proxyConn, err = t.ctl.GetProxy(); err != nil {
			log.Printf("Failded to get proxy connection: %v \n", err)
			return
		}
		//defer proxyConn.Close()
		log.Printf("Got proxy connectin %v \n", proxyConn.Id())

		startPxyMsg := &msg.StartProxy{
			Url:        t.url,
			ClientAddr: publicConn.RemoteAddr().String(),
		}
		if err = msg.WriteMsg(proxyConn, startPxyMsg); err != nil {
			log.Printf("Faild to write StartProxyMessage: %v, attempt %d\n", err, i)
			proxyConn.Close()
		} else {
			break
		}
	}
	if err != nil {
		log.Printf("Too many failures starting proxy connection")
		return
	}
	proxyConn.SetDeadline(time.Time{})
	_, _ = conn.Join(publicConn, proxyConn)

	log.Printf("join ok %v\n", startTime)
	proxyConn.Close()
}
