package main

import (
	"time"
	"net"
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

