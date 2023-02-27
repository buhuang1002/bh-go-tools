package bhsync

import "sync"

type keyMut struct {
	sync.Mutex
	int
}

type MapMutex struct {
	mut sync.Mutex
	m   map[interface{}]*keyMut
}

func NewMapMutex() *MapMutex {
	return &MapMutex{
		mut: sync.Mutex{},
		m:   map[interface{}]*keyMut{},
	}
}

func (m *MapMutex) Lock(k interface{}) {
	m.mut.Lock()
	v, ok := m.m[k]
	if !ok {
		v = &keyMut{sync.Mutex{}, 1}
		m.m[k] = v
		v.Lock()
		m.mut.Unlock()
		return
	}
	v.int++
	m.mut.Unlock()
	v.Lock()
	return
}

func (m *MapMutex) TryLock(k interface{}) bool {
	m.mut.Lock()
	v, ok := m.m[k]
	if !ok {
		v = &keyMut{sync.Mutex{}, 1}
		m.m[k] = v
		v.Lock()
		m.mut.Unlock()
		return true
	}

	m.mut.Unlock()
	return false
}

func (m *MapMutex) Unlock(k interface{}) {
	m.mut.Lock()
	defer m.mut.Unlock()
	v, ok := m.m[k]
	if !ok {
		panic("unreachable")
	}
	v.Unlock()
	v.int--
	if v.int == 0 {
		delete(m.m, k)
	}
}
