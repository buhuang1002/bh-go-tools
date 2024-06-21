package bhmap

import (
	"sync"

	"github.com/emirpasic/gods/v2/maps"
)

var _ maps.Map[string, any] = NewSyncMap[string, any](nil)

func NewSyncMap[K comparable, V any](m maps.Map[K, V]) *SyncMap[K, V] {
	return &SyncMap[K, V]{
		inner: m,
	}
}

type SyncMap[K comparable, V any] struct {
	inner maps.Map[K, V]
	m     sync.RWMutex
}

func (sm *SyncMap[K, V]) Put(key K, value V) {
	sm.m.Lock()
	defer sm.m.Unlock()

	sm.inner.Put(key, value)
}

func (sm *SyncMap[K, V]) Get(key K) (value V, found bool) {
	sm.m.RLock()
	defer sm.m.RUnlock()

	return sm.inner.Get(key)
}

func (sm *SyncMap[K, V]) Remove(key K) {
	sm.m.Lock()
	defer sm.m.Unlock()

	sm.inner.Remove(key)
}

func (sm *SyncMap[K, V]) Keys() []K {
	sm.m.RLock()
	defer sm.m.RUnlock()

	return sm.inner.Keys()
}

func (sm *SyncMap[K, V]) Empty() bool {
	sm.m.RLock()
	defer sm.m.RUnlock()

	return sm.inner.Empty()
}

func (sm *SyncMap[K, V]) Size() int {
	sm.m.RLock()
	defer sm.m.RUnlock()

	return sm.inner.Size()
}

func (sm *SyncMap[K, V]) Clear() {
	sm.m.Lock()
	defer sm.m.Unlock()

	sm.inner.Clear()
}

func (sm *SyncMap[K, V]) Values() []V {
	sm.m.RLock()
	defer sm.m.RUnlock()

	return sm.inner.Values()
}

func (sm *SyncMap[K, V]) String() string {
	sm.m.RLock()
	defer sm.m.RUnlock()

	return sm.inner.String()
}
