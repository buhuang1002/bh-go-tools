package bhsync

import "sync"

type keyRWMut struct {
	sync.RWMutex
	int
}

type MapRWMutex struct {
	mut sync.Mutex
	m   map[interface{}]*keyRWMut
}

func NewMapRWMutex() *MapRWMutex {
	return &MapRWMutex{
		mut: sync.Mutex{},
		m:   map[interface{}]*keyRWMut{},
	}
}

func (m *MapRWMutex) Lock(k interface{}) {
	m.mut.Lock()
	v, ok := m.m[k]
	if !ok {
		v = &keyRWMut{sync.RWMutex{}, 1}
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

func (m *MapRWMutex) RLock(k interface{}) {
	m.mut.Lock()
	v, ok := m.m[k]
	if !ok {
		v = &keyRWMut{sync.RWMutex{}, 1}
		m.m[k] = v
		v.RLock()
		m.mut.Unlock()
		return
	}

	v.int++
	m.mut.Unlock()
	v.RLock()
	return
}

func (m *MapRWMutex) TryLock(k interface{}) bool {
	m.mut.Lock()
	v, ok := m.m[k]
	if !ok {
		v = &keyRWMut{sync.RWMutex{}, 1}
		m.m[k] = v
		v.Lock()
		m.mut.Unlock()
		return true
	}

	m.mut.Unlock()
	return false
}

func (m *MapRWMutex) TryRLock(k interface{}) bool {
	m.mut.Lock()
	v, ok := m.m[k]
	if !ok {
		v = &keyRWMut{sync.RWMutex{}, 1}
		m.m[k] = v
		locked := v.TryRLock()
		if !locked {
			panic("unreachable")
		}

		m.mut.Unlock()
		return true
	}

	locked := v.TryRLock()
	if locked {
		v.int++
	}

	m.mut.Unlock()
	return locked
}

func (m *MapRWMutex) Unlock(k interface{}) {
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

func (m *MapRWMutex) RUnlock(k interface{}) {
	m.mut.Lock()
	defer m.mut.Unlock()
	v, ok := m.m[k]
	if !ok {
		panic("unreachable")
	}

	v.RUnlock()
	v.int--
	if v.int == 0 {
		delete(m.m, k)
	}
}
