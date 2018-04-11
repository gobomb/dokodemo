package main

import (
	"github.com/qiniu/log"
	"doko/conn"
	"doko/msg"
	"runtime"
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
	log.Printf("%v", config)

	run(config)

}

func run(config *Configuration) {
	model := newClientModel(config)
	model.run()

}

func newClientModel(config *Configuration) *ClientModel {
	return &ClientModel{
		serverAddr: config.ServerAddr,
		tunnelConfig:    config.Tunnels,
	}
}

type ClientModel struct {
	id           string
	//tunnels      map[string]Tunnel
	//updateStatus UpdateStatus
	//connStatus   ConnStatus
	//ctl          Controller
	serverAddr   string
	tunnelConfig map[string]*TunnelConfiguration
}

func (c *ClientModel) run() {
	//for {
		c.control()
	//}
}

func (c *ClientModel) control() {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("control recovering from failure %v", r)
		}
	}()
	var ctlConn conn.Conn

	ctlConn = conn.Dial(c.serverAddr)

	auth := &msg.Auth{
		ClientId:  c.id,
		OS:        runtime.GOOS,
		Arch:      runtime.GOARCH,
		//Version:   version.Proto,
		//MmVersion: version.MajorMinor(),
		//User:      c.authToken,
	}

	if err := msg.WriteMsg(ctlConn, auth); err != nil {
		log.Printf("[msg.WriteMsg] %v",err)
		panic(err)
	}

	//log.Println(n)
	//log.Println(err)

	//defer ctlConn.Close()
}
