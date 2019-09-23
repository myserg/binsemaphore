package binsemaphore

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestSemaphore(t *testing.T) {
	t.Run("Wait() for value/Post()", func(t *testing.T) {
		var s S
		var val int32
		exp := int32(1234)
		var wg sync.WaitGroup
		wg.Add(1)
		go func() {
			time.Sleep(10 * time.Millisecond)
			atomic.StoreInt32(&val, exp)
			s.Post()
			wg.Done()
		}()
		s.Wait()
		if atomic.LoadInt32(&val) != exp {
			t.Errorf("val %d expected %d", val, exp)
		}
		wg.Wait()
	})

	t.Run("WaitWithTimeout() for value/Post()", func(t *testing.T) {
		var s S
		var val int32
		exp := int32(12345)
		var wg sync.WaitGroup
		wg.Add(1)
		go func() {
			time.Sleep(10 * time.Millisecond)
			atomic.StoreInt32(&val, exp)
			s.Post()
			wg.Done()
		}()
		timedout := s.WaitWithTimeout(3 * time.Second)
		if timedout != false {
			t.Errorf("timedout is true")
		}
		if atomic.LoadInt32(&val) != exp {
			t.Errorf("val %d expected %d", val, exp)
		}
		wg.Wait()
	})

	t.Run("WaitWithTimeout() timeout", func(t *testing.T) {
		var s S
		timeout := 77 * time.Millisecond
		start := time.Now()
		timedout := s.WaitWithTimeout(timeout)
		elapsed := time.Since(start)
		if timedout != true {
			t.Errorf("timedout is not true")
		}
		if elapsed < timeout {
			t.Errorf("elapsed < timeout")
		}
	})
}
