package main

import (
	"github.com/qiniu/log"
	//"doko/conn"
)

type Options struct {
	hostname string
	protocol string
}

type Configuration struct {
	ServerAddr string                          `yaml:"server_addr,omitempty"`
	Tunnels    map[string]*TunnelConfiguration `yaml:"tunnels,omitempty"`
}

type TunnelConfiguration struct {
	Protocols  map[string]string `yaml:"proto,omitempty"`
	RemotePort uint16            `yaml:"remote_port,omitempty"`
}

func main() {
	protocols := make(map[string]string)
	protocols["tcp"] = "192.168.2.20:80"
	tunnels := make(map[string]*TunnelConfiguration)
	tunnels["szu"] = &TunnelConfiguration{
		Protocols:  protocols,
		RemotePort: 9999,
	}
	config := &Configuration{
		ServerAddr: "0.0.0.0:4443",
		Tunnels:    tunnels,
	}
	log.Printf("%v",config)
	control()
}

func control() {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("control recovering from failure %v", r)
		}
	}()
	//var ctlConn conn.Conn
}
