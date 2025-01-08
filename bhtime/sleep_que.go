package bhtime

import (
	"math"
	"sort"
	"time"
)

type SleepQue[T any] struct {
	l       []item[T]
	newCh   chan item[T]
	C       <-chan T
	closing chan struct{}
}

type item[T any] struct {
	gape time.Duration
	data T
}

func NewSleepQue[T any]() *SleepQue[T] {
	c := make(chan T)
	wq := &SleepQue[T]{
		l:       nil,
		newCh:   make(chan item[T]),
		closing: make(chan struct{}),
		C:       c,
	}

	go wq.startToWWW(c)
	return wq
}

func (q *SleepQue[T]) Put(d time.Duration, data T) {
	q.newCh <- item[T]{time.Duration(time.Now().UnixNano()) + d, data}

}

func (q *SleepQue[T]) startToWWW(c chan<- T) {
	var (
		toWait time.Duration
		data   T
	)

	reset := func() {
		var zeroData T
		toWait = time.Duration(math.MaxInt64)
		data = zeroData
	}

	set := func(item item[T]) {
		toWait = item.gape - time.Duration(time.Now().UnixNano())
		data = item.data
	}

	newOne := func(one item[T]) {
		q.l = append(q.l, one)
		sort.Slice(q.l, func(i, j int) bool {
			return q.l[i].gape < q.l[j].gape
		})

		set(q.l[0])
	}

	reset()

	_close := func() {
		reset()
		q.l = nil
		close(c)
		return
	}

	for {
		select {
		case <-q.closing:
			_close()
		case one := <-q.newCh:
			newOne(one)
		case <-time.After(toWait):
			select {
			case <-q.closing:
				_close()
			case one := <-q.newCh:
				newOne(one)

			case c <- data:
				q.l[0] = item[T]{}
				q.l = q.l[1:]
				if len(q.l) == 0 {
					reset()
				} else {
					set(q.l[0])
				}
			}
		}
	}

}

func (q *SleepQue[T]) Stop() {
	close(q.closing)
}
