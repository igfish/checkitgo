package checkit

import (
	"testing"
)

func TestGet(t *testing.T) {
	m := NewRBTreeMap[int, int]()
	m.Put(1, 2)
	m.Put(2, 3)
	if m.Len() != 2 {
		t.Fatalf("")
	}
	if v, ok := m.Get(1); !ok || v != 2 {
		t.Fatalf("")
	}
	if v, ok := m.Get(2); !ok || v != 3 {
		t.Fatalf("")
	}
}

func TestIter(t *testing.T) {
	m := NewRBTreeMap[int, int]()
	for i := 1000; i >= 0; i-- {
		m.Put(i, i*i)
	}
	itor := m.Iter()
	for i := 0; i <= 1000; i++ {
		pair, ok := itor.Next()
		if !ok || pair.K != i || pair.V != i*i {
			t.Fatalf("")
		}
	}
	if pair, ok := itor.Next(); pair != nil || ok {
		t.Fatalf("")
	}
}
