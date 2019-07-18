package pkg

import (
	"testing"
	"time"
)

func TestFutureGet(t *testing.T) {
	f := NewFuture()
	f.Submit(func() (interface{}, error) {
		time.Sleep(1 * time.Second)
		return "test", nil
	})

	v, _ := f.Get()
	if v != "test" {
		t.Errorf("Future Get retrun wrong %s", v)
	}
}

func TestFutureGetUntil(t *testing.T) {
	f := NewFuture()
	f.Submit(func() (interface{}, error) {
		time.Sleep(1 * time.Second)
		return "test", nil
	})

	_, err := f.GetUntil(100)
	if err == nil {
		t.Errorf("GetUntil timeout failed: %s", err)
	}

	time.Sleep(1 * time.Second)

	v, err := f.GetUntil(100)
	if v != "test" {
		t.Errorf("Future GetUntil failed: %s", v)
	}
}

func TestFutureState(t *testing.T) {
	f := NewFuture()

	if f.state != PENDING {
		t.Errorf("Future State error: %v", f.state)
	}

	stop := make(chan struct{})
	f.Submit(func() (interface{}, error) {
		for {
			select {
			case <-stop:
				return "test", nil
			default:
				time.Sleep(1 * time.Second)
			}
		}

	})

	if f.state != RUNNING {
		t.Errorf("Future State error: %v", f.state)
	}

	close(stop)
	f.Get()
	if f.state != FINISHED {
		t.Errorf("Future State error: %v", f.state)
	}
}
