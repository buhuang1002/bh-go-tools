package bhrunner

import (
	"context"
	"sync"
)

// GroupRunner manages multiple cancellable goroutines.
type GroupRunner struct {
	mu      sync.Mutex
	wg      sync.WaitGroup
	cancels []context.CancelFunc
}

// NewGroupRunner creates a new GroupRunner instance.
func NewGroupRunner() *GroupRunner {
	return &GroupRunner{}
}

// Go starts a new task with its own cancelable context.
func (g *GroupRunner) Go(parent context.Context, run func(ctx context.Context)) {
	ctx, cancel := context.WithCancel(parent)

	g.mu.Lock()
	g.cancels = append(g.cancels, cancel)
	g.wg.Add(1)
	g.mu.Unlock()

	go func() {
		defer g.wg.Done()
		run(ctx)
	}()
}

// Stop cancels all running tasks and waits for them to finish.
func (g *GroupRunner) Stop(ctx context.Context) error {
	g.mu.Lock()
	for _, cancel := range g.cancels {
		cancel()
	}
	g.cancels = nil
	g.mu.Unlock()

	done := make(chan struct{})
	go func() {
		g.wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}
