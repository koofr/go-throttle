package throttle_test

import (
	"sync"
	"testing"
	"time"

	"github.com/koofr/go-throttle"
)

func Test_NoPings(t *testing.T) {
	var wg sync.WaitGroup
	throttle := throttle.NewThrottle(10*time.Millisecond, false)

	count := 0

	wg.Add(1)
	go func() {
		defer wg.Done()
		for throttle.Next() {
			count += 1
		}
	}()

	time.Sleep(100 * time.Millisecond)
	throttle.Stop()

	wg.Wait()

	if count != 0 {
		t.Errorf("count = %v", count)
	}
}

func Test_MultiPingInOnePeriod(t *testing.T) {
	var wg sync.WaitGroup

	throttle := throttle.NewThrottle(10*time.Millisecond, false)
	count := 0

	wg.Add(1)
	go func() {
		defer wg.Done()
		for throttle.Next() {
			count += 1
		}
	}()

	for i := 0; i < 5; i++ {
		throttle.Trigger()
	}

	time.Sleep(50 * time.Millisecond)

	throttle.Stop()

	wg.Wait()

	if count != 1 {
		t.Errorf("count = %v", count)
	}
}

func Test_MultiPingInMultiplePeriod(t *testing.T) {
	var wg sync.WaitGroup

	throttle := throttle.NewThrottle(10*time.Millisecond, false)
	count := 0

	wg.Add(1)
	go func() {
		defer wg.Done()
		for throttle.Next() {
			count += 1
		}
	}()

	for i := 0; i < 5; i++ {
		time.Sleep(10 * time.Millisecond / 4)
		throttle.Trigger()
	}

	time.Sleep(50 * time.Millisecond)

	throttle.Stop()

	wg.Wait()

	if count != 2 {
		t.Errorf("count = %v", count)
	}
}

func Test_TrailingMultiPingInOnePeriod(t *testing.T) {
	var wg sync.WaitGroup

	throttle := throttle.NewThrottle(10*time.Millisecond, true)
	count := 0

	cond := sync.NewCond(&sync.Mutex{})

	wg.Add(1)
	go func() {
		defer wg.Done()
		for throttle.Next() {
			count += 1
			cond.Broadcast()
		}
	}()

	throttle.Trigger()

	cond.L.Lock()
	cond.Wait()
	throttle.Trigger()
	cond.L.Unlock()

	throttle.Trigger()
	throttle.Trigger()

	time.Sleep(50 * time.Millisecond)

	throttle.Stop()

	wg.Wait()

	if count != 2 {
		t.Errorf("count = %v", count)
	}
}

func Test_ThrottleFunc(t *testing.T) {
	count := 0

	throttle := throttle.ThrottleFunc(10*time.Millisecond, false, func() {
		count += 1
	})

	for i := 0; i < 5; i++ {
		throttle.Trigger()
	}

	time.Sleep(50 * time.Millisecond)

	throttle.Stop()

	if count != 1 {
		t.Errorf("count = %v", count)
	}
}
