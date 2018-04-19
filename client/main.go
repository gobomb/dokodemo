package main

import (
	"doko/conn"
	"doko/msg"
	"doko/util"
	"fmt"
	"github.com/qiniu/log"
	"runtime"
	"strings"
	"sync/atomic"
	"time"
)

const (
	defaultServerAddr   = "ngrokd.ngrok.com:443"
	defaultInspectAddr  = "127.0.0.1:4040"
	pingInterval        = 20 * time.Second
	maxPongLatency      = 15 * time.Second
	updateCheckInterval = 6 * time.Hour
	BadGateway          = `<html>
<body style="background-color: #97a8b9">
    <div style="margin:auto; width:400px;padding: 20px 60px; background-color: #D3D3D3; border: 5px solid maroon;">
        <h2>Tunnel %s unavailable</h2>
        <p>Unable to initiate connection to <strong>%s</strong>. A web server must be running on port <strong>%s</strong> to complete the tunnel.</p>
`
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
	//protocols["tcp"] = "192.168.2.20:80"
	protocols["tcp"] = "0.0.0.0:22"
	tunnels := make(map[string]*TunnelConfiguration)
	tunnels["ssh"] = &TunnelConfiguration{
		Protocols:  protocols,
		RemotePort: 9999,
	}
	config := &Configuration{
		ServerAddr: "0.0.0.0:4443",
		Tunnels:    tunnels,
	}

	run(config)

}

func run(config *Configuration) {
	model := newClientModel(config)
	model.run()

}

func newClientModel(config *Configuration) *ClientModel {
	return &ClientModel{
		serverAddr:   config.ServerAddr,
		tunnelConfig: config.Tunnels,
		tunnels:      make(map[string]Tunnel),
	}
}

type ClientModel struct {
	id      string
	tunnels map[string]Tunnel
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
	//defer func() {
	//	if r := recover(); r != nil {
	//		log.Printf("control recovering from failure %v", r)
	//	}
	//}()
	var ctlConn conn.Conn
	//var err error

	ctlConn = conn.Dial(c.serverAddr, "ctl")

	c.auth(ctlConn)

	c.reqTunnel(ctlConn)

	c.mainControlLoop(ctlConn)

}
func (c *ClientModel) proxy() {
	var (
		remoteConn conn.Conn
		err        error
	)

	remoteConn = conn.Dial(c.serverAddr, "pxy")
	if err != nil {
		log.Errorf("Failed to establish proxy connection: %v", err)
		return
	}
	defer remoteConn.Close()
	err = msg.WriteMsg(remoteConn, &msg.RegProxy{ClientId: c.id})
	if err != nil {
		log.Errorf("Failed to write RegProxy: %v", err)
		return
	}

	var startPxy msg.StartProxy
	if err = msg.ReadMsgInto(remoteConn, &startPxy); err != nil {
		log.Errorf("Server failed to write StartProxy: %v", err)
		return
	}
	tunnel, ok := c.tunnels[startPxy.Url]
	if !ok {
		log.Errorf("Couldn't find tunnel for proxy: %s", startPxy.Url)
		return
	}

	//start := time.Now()
	localConn := conn.Dial(tunnel.LocalAddr, "prv")
	defer localConn.Close()
	//locConn = wrapcon
	conn.Join(localConn, remoteConn)
}

func (c *ClientModel) mainControlLoop(ctlConn conn.Conn) {
	var err error
	for {
		var rawMsg msg.Message
		if rawMsg, err = msg.ReadMsg(ctlConn); err != nil {
			panic(err)
		}

		switch m := rawMsg.(type) {
		case *msg.ReqProxy:
			go c.proxy()
		case *msg.Pong:
			atomic.StoreInt64(&lastPong, time.Now().UnixNano())
		case *msg.NewTunnel:
			if m.Error != "" {
				emsg := fmt.Sprintf("Server failed to allocate tunnel: %s", m.Error)
				log.Error(emsg)
				//c.ctl.Shutdown(emsg)
				continue
			}

			tunnel := Tunnel{
				PublicUrl: m.Url,
				LocalAddr: reqIdToTunnelConfig[m.ReqId].Protocols[m.Protocol],
			}
			c.tunnels[tunnel.PublicUrl] = tunnel
			log.Info(tunnel)
			log.Printf("Tunnel established at %v", tunnel.PublicUrl)
		default:
			log.Warnf("Ignoring unknown control message %v ", m)
		}
	}
}

type Tunnel struct {
	PublicUrl string
	//Protocol  proto.Protocol
	LocalAddr string
}

var reqIdToTunnelConfig map[string]*TunnelConfiguration

func (c *ClientModel) reqTunnel(ctlConn conn.Conn) {
	reqIdToTunnelConfig = make(map[string]*TunnelConfiguration)

	for _, config := range c.tunnelConfig {
		log.Info(config)
		var protocols []string
		for proto, _ := range config.Protocols {
			protocols = append(protocols, proto)
		}
		reqId := util.GenerateRandomString()
		reqTunnel := &msg.ReqTunnel{
			ReqId:      reqId,
			Protocol:   strings.Join(protocols, "+"),
			RemotePort: config.RemotePort,
		}

		if err := msg.WriteMsg(ctlConn, reqTunnel); err != nil {
			panic(err)
		}
		reqIdToTunnelConfig[reqTunnel.ReqId] = config
	}

	lastPong = time.Now().UnixNano()

	go c.heartbeat(&lastPong, ctlConn)

}

var lastPong int64

// Hearbeating to ensure our connection ngrokd is still live
func (c *ClientModel) heartbeat(lastPongAddr *int64, conn conn.Conn) {
	lastPing := time.Unix(atomic.LoadInt64(lastPongAddr)-1, 0)
	ping := time.NewTicker(pingInterval)
	pongCheck := time.NewTicker(time.Second)

	defer func() {
		conn.Close()
		ping.Stop()
		pongCheck.Stop()
	}()

	for {
		select {
		case <-pongCheck.C:
			lastPong := time.Unix(0, atomic.LoadInt64(lastPongAddr))
			needPong := lastPong.Sub(lastPing) < 0
			pongLatency := time.Since(lastPing)

			if needPong && pongLatency > maxPongLatency {
				log.Printf("Last ping: %v, Last pong: %v", lastPing, lastPong)
				log.Printf("Connection stale, haven't gotten PongMsg in %d seconds", int(pongLatency.Seconds()))
				return
			}

		case <-ping.C:
			err := msg.WriteMsg(conn, &msg.Ping{})
			if err != nil {
				log.Printf("Got error %v when writing PingMsg", err)
				return
			}
			lastPing = time.Now()
		}
	}
}

func (c *ClientModel) auth(ctlConn conn.Conn) {
	var err error
	auth := &msg.Auth{
		OS:   runtime.GOOS,
		Arch: runtime.GOARCH,
		//Version:   version.Proto,
		//MmVersion: version.MajorMinor(),
		//User:      c.authToken,
	}

	if err = msg.WriteMsg(ctlConn, auth); err != nil {
		log.Printf("[msg.WriteMsg] %v", err)
		panic(err)
	}

	var authResp msg.AuthResp
	if err = msg.ReadMsgInto(ctlConn, &authResp); err != nil {
		log.Printf("[msg.ReadMsgInto] %v", err)
		panic(err)
	}
	c.id = authResp.ClientId
	if authResp.Error != "" {
		emsg := fmt.Sprintf("Faild to authenticate to server: %s", authResp.Error)
		log.Printf("[] %v", emsg)
		return
	}
}
