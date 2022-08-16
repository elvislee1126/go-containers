package rollingwindow

import (
	linkedmap "github.com/elvislee1126/gocontainers/linked-map"
)

type windowSlot[K comparable, V any] struct {
	paris   *linkedmap.LinkedMap[K, V]
	slotIdx int64
}

func newWindowSlot[K comparable, V any](slotIdx int64) *windowSlot[K, V] {
	lm := linkedmap.New[K, V]()
	slot := &windowSlot[K, V]{
		paris:   lm,
		slotIdx: slotIdx,
	}
	return slot
}

func (slot *windowSlot[K, V]) Store(k K, v V) {
	slot.paris.Store(k, v)
}

func (slot *windowSlot[K, V]) Load(k K) (V, bool) {
	return slot.paris.Load(k)
}

func (slot *windowSlot[K, V]) Size() int {
	return slot.paris.Len()
}
