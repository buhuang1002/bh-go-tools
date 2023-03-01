package bhsync

import "sync"

type keyMut struct {
	sync.Mutex
	int
}

type MapMutex[T comparable] struct {
	mut sync.Mutex
	m   map[T]*keyMut
}

func NewMapMutex[T comparable]() *MapMutex[T] {
	return &MapMutex[T]{
		mut: sync.Mutex{},
		m:   map[T]*keyMut{},
	}
}

func (m *MapMutex[T]) Lock(k T) {
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

func (m *MapMutex[T]) TryLock(k T) bool {
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

func (m *MapMutex[T]) Unlock(k T) {
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
