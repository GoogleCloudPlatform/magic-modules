package tpgresource

import (
	"sync"
)

// WorkerPool runs submitted functions concurrently up to a fixed worker limit.
type WorkerPool struct {
	sem chan struct{}
	wg  sync.WaitGroup
}

// NewWorkerPool creates a WorkerPool with at most maxWorkers concurrent workers.
// If maxWorkers < 1, it defaults to 1.
func NewWorkerPool(maxWorkers int) *WorkerPool {
	if maxWorkers < 1 {
		maxWorkers = 1
	}
	return &WorkerPool{
		sem: make(chan struct{}, maxWorkers),
	}
}

// Submit schedules fn for execution in the worker pool.
// If fn is nil, Submit is a no-op.
func (wp *WorkerPool) Submit(fn func()) {
	if fn == nil {
		return
	}
	wp.wg.Add(1)
	go func() {
		defer wp.wg.Done()
		wp.sem <- struct{}{}
		defer func() { <-wp.sem }()
		fn()
	}()
}

// StopWait waits for all submitted functions to complete.
func (wp *WorkerPool) StopWait() {
	wp.wg.Wait()
}
