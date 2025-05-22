package host

import (
	"errors"
	"sync"

	"infra-scheduler/pkg/config"
)

type Host struct {
	ID      string
	Address string
	CPU     int
	Memory  int
	Storage int
	mu      sync.Mutex
}

type Manager struct {
	hosts []*Host
}

// NewManager creates a Manager with initial capacities
func NewManager(cfgs []config.HostConfig) (*Manager, error) {
	if len(cfgs) == 0 {
		return nil, errors.New("no hosts configured")
	}
	mgr := &Manager{}
	for _, hc := range cfgs {
		mgr.hosts = append(mgr.hosts, &Host{
			ID:      hc.ID,
			Address: hc.Address,
			CPU:     hc.CPU,
			Memory:  hc.Memory,
			Storage: hc.Storage,
		})
	}
	return mgr, nil
}

// Hosts returns the list of hosts
func (m *Manager) Hosts() []*Host {
	return m.hosts
}

// Reserve attempts to allocate resources on the host
func (h *Host) Reserve(cores, memory, disk int) bool {
	h.mu.Lock()
	defer h.mu.Unlock()
	if h.CPU >= cores && h.Memory >= memory && h.Storage >= disk {
		h.CPU -= cores
		h.Memory -= memory
		h.Storage -= disk
		return true
	}
	return false
}
