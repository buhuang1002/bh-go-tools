package bhsync

import "sync"

type keyRWMut struct {
	sync.RWMutex
	int
}

type MapRWMutex[T comparable] struct {
	mut sync.Mutex
	m   map[T]*keyRWMut
}

func NewMapRWMutex[T comparable]() *MapRWMutex[T] {
	return &MapRWMutex[T]{
		mut: sync.Mutex{},
		m:   map[T]*keyRWMut{},
	}
}

func (m *MapRWMutex[T]) Lock(k T) {
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

func (m *MapRWMutex[T]) RLock(k T) {
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

func (m *MapRWMutex[T]) TryLock(k T) bool {
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

func (m *MapRWMutex[T]) TryRLock(k T) bool {
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

func (m *MapRWMutex[T]) Unlock(k T) {
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

func (m *MapRWMutex[T]) RUnlock(k T) {
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
