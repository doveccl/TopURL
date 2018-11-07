package main

type PSI struct {
	s string
	i int
}

type PSIHeap []PSI

func (h PSIHeap) Len() int {
	return len(h)
}
func (h PSIHeap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}
func (h PSIHeap) Less(i, j int) bool {
	return h[i].i < h[j].i
}

func (h *PSIHeap) Push(i interface{}) {
	*h = append(*h, i.(PSI))
}
func (h *PSIHeap) Pop() interface{} {
	n := len(*h)
	x := (*h)[n - 1]
	*h = (*h)[:n - 1]
	return x
}
