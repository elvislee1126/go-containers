package rollingwindow

import (
	"fmt"
	"math/rand"
	"sync"
	"testing"
	"time"
)

func TestRollingWindow(t *testing.T) {
	window := New[string, int](
		WithSlotAmountAndSize(6, time.Second),
		WithTimeProvider(func() time.Time {
			return time.Now()
		}),
	)

	window.Set("elvis", 20)

	time.Sleep(time.Second * 2)
	val, ok := window.Get("elvis", WithCurrentWindow())
	if ok {
		t.Logf("get value success, value is %v", val)
	} else {
		t.Fatalf("got value failed, value is %v", val)
	}

	time.Sleep(time.Second * 5)
	if _, ok = window.Get("elvis", WithCurrentWindow()); ok {
		t.Fatalf("valeu should be drain, but still got it")
	} else {
		t.Logf("key is successfuly drained")
	}
}

func TestGetAfter(t *testing.T) {
	window := New[string, int](
		WithSlotAmountAndSize(10, time.Second),
		WithTimeProvider(func() time.Time {
			return time.Now()
		}),
		WithVerbose(),
	)
	window.Set("elvis", 20)
	time.Sleep(time.Second * 3)
	val, ok := window.Get("elvis", WithCurrentWindow())
	if ok {
		t.Logf("get value success, value is %d", val)
	} else {
		t.Fatalf("get value failed, value is %d", val)
	}

	now := time.Now()
	dur := time.Second * -2
	val, ok = window.Get("elvis", WithAfter(now.Add(dur)))
	if ok {
		t.Fatalf("still get value, value is %d", val)
	} else {
		t.Logf("not get value")
	}

	time.Sleep(time.Second * 8)
	val, ok = window.Get("elvis", WithAfter(now.Add(dur)))
	if ok {
		t.Fatalf("still get value, value is %d", val)
	} else {
		t.Logf("not get value")
	}

}

func TestSetNX(t *testing.T) {
	window := New[string, int](
		WithSlotAmountAndSize(10, time.Second),
		WithTimeProvider(func() time.Time {
			return time.Now()
		}),
		WithVerbose(),
	)
	wg := sync.WaitGroup{}
	wg.Add(100)
	okCounter := 0
	for i := 0; i < 100; i++ {
		go func() {
			ok, _ := window.Set("elvis", 20, WithSetNX(true))
			if ok {
				okCounter++
			}
			wg.Done()
		}()
	}
	wg.Wait()
	if okCounter != 1 {
		t.Fatalf("期望值是只有一个线程 set 成功，但是实际上有 %d", okCounter)
	} else {
		t.Logf("%d", okCounter)
	}
}

func TestDrained(t *testing.T) {
	dur := time.Second * 5
	window := New[string, int](
		WithSlotAmountAndSize(1, dur),
		WithTimeProvider(func() time.Time {
			return time.Now()
		}),
	)
	window.Set("elvis", 20)

	val, ok := window.Get("elvis", WithCurrentWindow())
	if !ok || val != 20 {
		t.FailNow()
		return
	}

	time.Sleep(dur)
	val, ok = window.Get("elvis", WithCurrentWindow())
	if ok || val == 20 {
		t.Logf("window size is %s, after %s, but still got value, value is %v", dur.String(), dur.String(), val)
		t.FailNow()
		return
	}
}

func BenchmarkRollingWindow(b *testing.B) {
	window := New[string, int](
		WithSlotAmountAndSize(5, time.Millisecond*500),
		WithTimeProvider(func() time.Time {
			return time.Now()
		}),
	)
	b.RunParallel(func(p *testing.PB) {
		for p.Next() {
			num := rand.Int()
			key := fmt.Sprintf("key-%d", num)
			window.Set(key, num)
		}
		b.Logf("size: %d", window.slots.Len())
	})
}
