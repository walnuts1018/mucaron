package safemutex

import "sync"

type Mutex struct {
	mu     sync.Mutex
	locked bool
}

func NewMutex() *Mutex {
	return &Mutex{}
}

func (m *Mutex) Lock() {
	m.mu.Lock()
	m.locked = true
}

func (m *Mutex) Unlock() {
	if m.locked {
		m.locked = false
		m.mu.Unlock()
	} else {
		return
	}
}
