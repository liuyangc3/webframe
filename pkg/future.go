package pkg

import (
	"errors"
	"time"
)

type futureState int

// The  future states
const (
	PENDING futureState = iota
	RUNNING
	CANCELLED
	FINISHED
)

type futureResult struct {
	val interface{}
	err error
}

// Future represents a value (or error) to be
// available at some later time.
type Future struct {
	// (result, err) represents future result
	result interface{}
	err    error

	state futureState

	done   chan struct{}
	cancel chan struct{}

	// OnComplete callback, called when the future is cancelled
	// or finishes running.
	callback func()
}

// NewFuture creates Future
func NewFuture() *Future {

	future := Future{
		state:    PENDING,
		done:     make(chan struct{}),
		cancel:   make(chan struct{}),
		callback: nil,
		result:   futureResult{},
	}
	return &future
}

// Submit a function and run it
func (f *Future) Submit(run func() (interface{}, error)) {
	f.state = RUNNING

	go func() {
		defer func() {
			close(f.done)
			f.invokeCallback()
		}()

		select {
		case <-f.cancel:
			return
		default:
			f.result, f.err = run()
			f.state = FINISHED
		}
	}()
}

// Get the result of Future
func (f *Future) Get() (interface{}, error) {
	select {
	case <-f.done:
		return f.result, f.err
	case <-f.cancel:
		return nil, errors.New("canceled")
	}
}

// GetUntil result of Future return timeout error after d ms
func (f *Future) GetUntil(d time.Duration) (interface{}, error) {
	select {
	case <-time.After(d * time.Millisecond):
		return nil, errors.New("timeout")
	case <-f.done:
		return f.result, f.err
	case <-f.cancel:
		return nil, errors.New("canceled")
	}
}

// Add callback funcion to the future, callback will be
// called when the future is cancelled or finishes running.
func (f *Future) onComplete(callback func()) {
	f.callback = callback
}

func (f *Future) invokeCallback() {
	if f.callback != nil {
		f.callback()
	}
}

// Cancel the future if possible.
// Returns true if the future was cancelled, false otherwise. A future
// cannot be cancelled if it is running or has already completed.
func (f *Future) Cancel() bool {
	defer f.invokeCallback()

	if f.state == RUNNING || f.state == FINISHED {
		return false
	}
	if f.state == CANCELLED {
		return true
	}

	close(f.cancel)
	f.state = CANCELLED

	return true
}

// IsCancelled returns true if the call was successfully cancelled.
func (f *Future) IsCancelled() bool {
	return f.state == CANCELLED
}

// IsDone returns true if the call was successfully cancelled or finished.
func (f *Future) IsDone() bool {
	return f.state == FINISHED || f.IsCancelled()
}
