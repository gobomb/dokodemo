package server

import (
	"sync"
	"fmt"
	"log"
)

type TunnelRegistry struct {
	tunnels map[string]*Tunnel
	sync.RWMutex
}

func (r *TunnelRegistry) Register(url string, t *Tunnel) error {
	r.Lock()
	defer r.Unlock()

	if r.tunnels[url] != nil {
		return fmt.Errorf("the tunnel %s is already registered", url)
	}

	r.tunnels[url] = t
	return nil
}

type ControlRegistry struct {
	controls map[string]*Control
	sync.RWMutex
}

func NewControlRegistry() *ControlRegistry {
	return &ControlRegistry{
		controls: make(map[string]*Control),
	}
}

func (r *ControlRegistry) Add(clientId string, ctl *Control) (oldCtl *Control){
	r.Lock()
	defer r.Unlock()

	oldCtl = r.controls[clientId]
	if oldCtl !=nil{
		oldCtl.Replaced(ctl)
	}

	r.controls[clientId] =ctl
	log.Printf("Registered control with id %s", clientId)
	return
}

