package main

import (
	"container/heap"
	"math/rand"
	"testing"
	"time"
)

func TestHeap(t *testing.T) {
	s := rand.NewSource(time.Now().Unix())
	r := rand.New(s)
	h := &PSIHeap{}
	heap.Init(h)
	for i := 0; i < 1e5; i++ {
		heap.Push(h, PSI{"str", r.Intn(1e7)})
		if h.Len() != i + 1 {
			t.Error("length error")
		}
	}
	prev := -1
	for i := 0; i < 1e5; i++ {
		p := heap.Pop(h).(PSI)
		if p.i < prev {
			t.Error("sort error")
		}
		prev = p.i
	}
}
