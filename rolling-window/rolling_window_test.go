package rollingwindow

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
)

func TestRollingWindow(t *testing.T) {
	window := NewRollingWindow[string, int](
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

func BenchmarkRollingWindow(b *testing.B) {
	window := NewRollingWindow[string, int](
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
