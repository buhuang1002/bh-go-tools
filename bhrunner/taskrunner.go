package bhrunner

import (
	"context"
	"errors"
	"sync"
)

// TaskRunner manages a single cancellable goroutine.
type TaskRunner struct {
	mu     sync.Mutex
	cancel context.CancelFunc
	done   chan struct{}
}

// NewTaskRunner creates a new TaskRunner instance.
func NewTaskRunner() *TaskRunner {
	return &TaskRunner{}
}

// Go starts the task. If already running, returns an error.
func (t *TaskRunner) Go(parent context.Context, run func(ctx context.Context)) error {
	t.mu.Lock()
	if t.cancel != nil {
		t.mu.Unlock()
		return errors.New("task already running")
	}

	ctx, cancel := context.WithCancel(parent)
	t.cancel = cancel
	t.done = make(chan struct{})
	t.mu.Unlock()

	go func() {
		run(ctx)

		t.mu.Lock()
		close(t.done)
		t.cancel = nil
		t.done = nil
		t.mu.Unlock()
	}()

	return nil
}

// Stop cancels the task and waits for it to finish.
func (t *TaskRunner) Stop(ctx context.Context) error {
	t.mu.Lock()
	if t.cancel == nil {
		t.mu.Unlock()
		return nil
	}
	cancel := t.cancel
	done := t.done
	t.mu.Unlock()

	cancel()

	select {
	case <-done:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}
