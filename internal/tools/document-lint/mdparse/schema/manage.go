package schema

import "sync"

// manage multiple packages

var manager = NewManager()

type Manager struct {
	Packs sync.Map // import path to Package
	mux   sync.RWMutex
}

func NewManager() *Manager {
	return &Manager{
		Packs: sync.Map{},
	}
}

// get package for import path
func (m *Manager) GetPackForPath(importPath string) *Package {
	if v, ok := m.Packs.Load(importPath); ok {
		return v.(*Package)
	}
	// one import path only init one time
	pack := NewPackage(importPath)
	m.Packs.Store(importPath, pack)
	return pack
}
