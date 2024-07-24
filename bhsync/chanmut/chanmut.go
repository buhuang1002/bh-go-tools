package chanmut

import "sync"

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

type chanMutexState struct {
	locked      chan struct{}
	interrupted chan struct{}
	toInterrupt chan struct{}
	m           chan struct{}
	once        sync.Once
}

func (cls *chanMutexState) Reset() bool {
	cls.once.Do(func() {
		close(cls.toInterrupt)
	})

	select {
	case <-cls.interrupted:
		return true
	case <-cls.locked:
		return false
	}
}

func (cls *chanMutexState) Done() <-chan struct{} {
	return cls.locked
}

func (cls *chanMutexState) IsLocked() bool {
	select {
	case <-cls.locked:
		return true
	default:
		return false
	}
}

func (cls *chanMutexState) Unlock() {
	<-cls.m
}

type ChanMutex interface {
	Lock() ChanMutexState
}

type ChanMutexState interface {
	Reset() bool
	Done() <-chan struct{}
	Unlock()
}
