package pkg

import (
	"testing"
	"time"
)

func TestFuture(t *testing.T) {
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
