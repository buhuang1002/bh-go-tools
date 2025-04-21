package mergecaller

import (
	"cmp"

	"github.com/buhuang1002/bh-go-tools/bhmap"
	"github.com/buhuang1002/bh-go-tools/bhsync"
	"github.com/emirpasic/gods/v2/maps/treemap"
)

func NewMergeCaller[K cmp.Ordered, V any]() *MergeCaller[K, V] {
	return &MergeCaller[K, V]{
		results: bhmap.NewSyncMap[K, *mergeResult[V]](treemap.New[K, *mergeResult[V]]()),
		m:       bhsync.NewMapMutex[K](),
	}
}

type MergeCaller[K cmp.Ordered, V any] struct {
	results *bhmap.SyncMap[K, *mergeResult[V]]
	m       *bhsync.MapMutex[K]
}

func (mi *MergeCaller[K, V]) Call(key K, f func() (V, error)) (V, error) {
	mi.m.Lock(key)
	result, ok := mi.results.Get(key)
	if ok {
		mi.m.Unlock(key)
		return result.waitReturn()
	}

	result = new_mergeResult[V]()
	mi.results.Put(key, result)
	mi.m.Unlock(key)

	result.result(f())
	mi.results.Remove(key)

	return result.waitReturn()
}

func new_mergeResult[V any]() *mergeResult[V] {
	return &mergeResult[V]{
		finished: make(chan struct{}),
	}
}

type mergeResult[V any] struct {
	v        V
	err      error
	finished chan struct{}
}

func (mr *mergeResult[V]) result(v V, err error) {
	mr.v = v
	mr.err = err
	close(mr.finished)
}

func (mr *mergeResult[V]) waitReturn() (v V, err error) {
	select {
	case <-mr.finished:
	}

	return mr.v, mr.err
}
