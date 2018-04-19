package server

import (
	"fmt"
	"log"
	"sync"
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

func (r *TunnelRegistry) Del(url string) {
	r.Lock()
	defer r.Unlock()
	delete(r.tunnels, url)
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

func (r *ControlRegistry) Add(clientId string, ctl *Control) (oldCtl *Control) {
	r.Lock()
	defer r.Unlock()

	oldCtl = r.controls[clientId]
	if oldCtl != nil {
		oldCtl.Replaced(ctl)
	}

	r.controls[clientId] = ctl
	log.Printf("Registered control with id %s", clientId)
	return
}

func (r *ControlRegistry) Get(clientId string) *Control {
	r.RLock()
	defer r.RUnlock()
	return r.controls[clientId]
}

func (r *ControlRegistry) Del(clientId string) error {
	r.Lock()
	defer r.Unlock()
	if r.controls[clientId] == nil {
		return fmt.Errorf("no control found for client id: %s", clientId)
	} else {
		log.Printf("Removed control registry id %s", clientId)
		delete(r.controls, clientId)
		return nil
	}
}
