package server

import (
	"time"
	"github.com/qiniu/log"
)

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
	IP         string
	Status     bool
	Tuns       []string
	Ctls       []string
	TunnelAddr string
	Domain     string
	Consnum    int
}

func GetInfo() Info {
	if !S.StatusOn {
		return Info{Status: false}
	}
	var (
		tuns []string
		ctls []string
	)
	for k, _ := range tunnelRegistry.tunnels {
		log.Info(k)
		tuns = append(tuns, k)
	}
	for k, _ := range controlRegistry.controls {
		log.Info(k)
		ctls = append(ctls, k)
	}
	ls := 0
	for k, _ := range listeners {
		log.Infof("listener:%v", k)
		ls++
	}
	return Info{
		IP:         "120.78.62.175",
		Status:     true,
		Tuns:       tuns,
		Ctls:       ctls,
		TunnelAddr: opts.tunnelAddr,
		Domain:     opts.domain,
		Consnum:    ls,
	}
}
