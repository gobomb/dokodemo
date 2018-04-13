package server

import (
	"time"
	"net"
	"log"
	"sync/atomic"
)

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

func (t *Tunnel)Shutdown(){
	log.Printf("Shutting down")
	atomic.StoreInt32(&t.closing,1)
	if t.listener!=nil{
		t.listener.Close()
	}
	tunnelRegistry.Del(t.url)
}