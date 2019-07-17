package pkg

import (
	"errors"
	"time"
)

type futureResult struct {
	get func() (interface{}, error)
}

// Future represents a value (or error) to be
// available at some later time.
type Future struct {
	cancelCh chan struct{}

	// Future is done
	done bool

	// OnComplete callback, called when the future is cancelled
	// or finishes running.
	callback func()

	result futureResult
}

// NewFuture creates Future
func NewFuture() *Future {
	return &Future{
		cancelCh: make(chan struct{}, 1),
		done:     false,
		callback: nil,
		result:   futureResult{},
	}
}

// Submit a function and run it
func (f *Future) Submit(callable func() (interface{}, error)) {
	var result interface{}
	var err error

	c := make(chan struct{}, 1)

	go func() {
		defer func() {
			close(c)
			f.done = true
			if f.callback != nil {
				// onComplete callback
				f.callback()
			}
		}()

		result, err = callable()
	}()

	f.result = futureResult{
		get: func() (interface{}, error) {
			select {
			case <-c:
				return result, err
			}
		},
	}
}

// Get result of Future
func (f *Future) Get() (interface{}, error) {
	return f.result.get()
}

// GetTimeout result of Future with timeout
func (f *Future) GetTimeout(timeout time.Duration) (interface{}, error) {
	select {
	case <-time.After(timeout):
		return nil, errors.New("timeout")
	default:
		return f.result.get()
	}
}

// Add callback funcion to the future, callback will be
// called when the future is cancelled or finishes running.
func (f *Future) onComplete(callback func()) {
	f.callback = callback
}

// Cancel cancels a future, can not be canceled when future is done.
func (f *Future) Cancel() error {
	if f.done {
		return errors.New("Cannot cancel future when status is Done")
	}
	return nil
}

// IsCancelled returns true if the call was successfully cancelled.
func (f *Future) IsCancelled() bool {
	// TODO
}

// IsDone returns true if the call was successfully cancelled
func (f *Future) IsDone() bool {
	return f.done
}
