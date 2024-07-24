package chanmut

import (
	"cmp"
	"sync"

	"github.com/emirpasic/gods/v2/maps/treemap"
)

func NewMapChanMutex[T cmp.Ordered]() *MapChanMutex[T] {
	return &MapChanMutex[T]{
		mut: sync.Mutex{},
		m:   treemap.New[T, *keyMut](),
	}
}

type keyMut struct {
	ChanMutex
	int
}

type MapChanMutex[T cmp.Ordered] struct {
	mut sync.Mutex
	m   *treemap.Map[T, *keyMut]
}

func (m *MapChanMutex[T]) Lock(k T) ChanMutexState {
	m.mut.Lock()

	v, ok := m.m.Get(k)
	if !ok {
		v = &keyMut{NewChanLock(), 1}
		m.m.Put(k, v)

		chstat := v.Lock()
		<-chstat.Done()

		m.mut.Unlock()
		return m.newMapChanMutexState(k, chstat)
	}

	v.int++
	m.mut.Unlock()
	return m.newMapChanMutexState(k, v.Lock())
}

func (m *MapChanMutex[T]) TryLock(k T) (ChanMutexState, bool) {
	m.mut.Lock()
	v, ok := m.m.Get(k)

	if !ok {
		v = &keyMut{NewChanLock(), 1}
		m.m.Put(k, v)

		chstat := v.Lock()
		<-chstat.Done()

		m.mut.Unlock()
		return m.newMapChanMutexState(k, chstat), true
	}

	m.mut.Unlock()
	return nil, false
}

func (m *MapChanMutex[T]) decrease(k T) {
	m.mut.Lock()
	defer m.mut.Unlock()

	v, ok := m.m.Get(k)
	if !ok {
		panic("unreachable")
	}

	v.int--
	if v.int == 0 {
		m.m.Remove(k)
	}
}

func (m *MapChanMutex[T]) newMapChanMutexState(k T, cmState ChanMutexState) ChanMutexState {
	return &mapChanMutexState[T]{
		key:     k,
		once:    sync.Once{},
		m:       m,
		cmState: cmState,
	}
}

type mapChanMutexState[T cmp.Ordered] struct {
	key     T
	once    sync.Once
	m       *MapChanMutex[T]
	cmState ChanMutexState
}

func (ms *mapChanMutexState[T]) Reset() bool {
	ok := ms.cmState.Reset()
	if ok {
		ms.once.Do(func() {
			ms.m.decrease(ms.key)
		})
	}

	return ok
}

func (ms *mapChanMutexState[T]) Done() <-chan struct{} {
	return ms.cmState.Done()
}

func (ms *mapChanMutexState[T]) Unlock() {
	ms.cmState.Unlock()
	ms.m.decrease(ms.key)
}

var _ ChanMutexState = (*mapChanMutexState[int])(nil)
