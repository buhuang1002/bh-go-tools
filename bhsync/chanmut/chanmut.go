package chanmut

import "sync"

type ChanMutex interface {
	Lock() ChanMutexState
}

func NewChanLock() ChanMutex {
	return &chanMutex{
		m: make(chan struct{}, 1),
	}
}

type chanMutex struct {
	m chan struct{}
}

func (m *chanMutex) Lock() ChanMutexState {
	locked := make(chan struct{})
	interrupted := make(chan struct{})
	toInterrupt := make(chan struct{})

	state := &chanMutexState{
		locked:      locked,
		interrupted: interrupted,
		toInterrupt: toInterrupt,
		m:           m.m,
	}

	go func() {
		select {
		case m.m <- struct{}{}:
			close(locked)
		case <-toInterrupt:
			close(interrupted)
		}
	}()

	return state
}

type ChanMutexState interface {
	Reset()
	Done() <-chan struct{}
}

type chanMutexState struct {
	locked      chan struct{}
	interrupted chan struct{}
	toInterrupt chan struct{}
	m           chan struct{}
	once        sync.Once
}

func (cls *chanMutexState) Reset() {
	cls.once.Do(func() {
		close(cls.toInterrupt)
		select {
		case <-cls.locked:
			<-cls.m
		case <-cls.interrupted:
		}
	})
}

func (cls *chanMutexState) Done() <-chan struct{} {
	return cls.locked
}
