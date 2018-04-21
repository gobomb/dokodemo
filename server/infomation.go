package server

import "time"

type InfoTunnel struct {
	// time when the tunnel was opened
	start time.Time

	// 公网 url
	Url string

	// tcp 监听
	Listener string

	// closing
	Closing int32
}

type Info struct {
	Tuns       []string
	Ctls       []string
	TunnelAddr string
	Domain     string
	consnum    int
}
