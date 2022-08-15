package linkedmap

import "testing"

func TestLinkedMapSetAndGet(t *testing.T) {
	lm := New[string, int]()
	lm.Store("elvis", 18)
	lm.Store("zoe", 20)
	val, _ := lm.Load("elvis")
	if val != 18 {
		t.FailNow()
	}
}

func TestLinkedMapSetDelete(t *testing.T) {
	lm := New[string, int]()
	lm.Store("elvis", 18)
	val, _ := lm.Delete("elvis")
	if val != 18 {
		t.FailNow()
	}
	_, ok := lm.Load("elvis")
	if ok {
		t.FailNow()
	}
}

func TestLinkedMapOrdered(t *testing.T) {
	type Entry struct {
		Key   string
		Value int
		Index int
	}
	enteries := map[string]Entry{
		"elvis": {"elvis", 18, 2},
		"zoe":   {"zoe", 20, 1},
		"amy":   {"amy", 21, 0},
	}
	lm := New[string, int]()
	for _, entry := range enteries {
		lm.Store(entry.Key, entry.Value)
	}
	keys := lm.Keys()
	for idx, k := range keys {
		sample := enteries[k]
		if sample.Index != idx {
			t.Errorf("key '%s' index should be %d, but got %d", k, sample.Index, idx)
		}
	}
}
