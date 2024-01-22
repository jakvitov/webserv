package cache

import "container/heap"

// Heap implementation for rebalancing the cache
type CachedFileHeap []*CachedFile

func (h CachedFileHeap) Len() int           { return len(h) }
func (h CachedFileHeap) Less(i, j int) bool { return h[i].GetCoefficient() < h[j].GetCoefficient() }
func (h CachedFileHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *CachedFileHeap) Push(x interface{}) {
	*h = append(*h, x.(*CachedFile))
	heap.Fix(h, len(*h)-1)
}

func (h *CachedFileHeap) Pop() interface{} {
	n := len(*h) - 1
	x := (*h)[n]
	*h = (*h)[:n]
	heap.Fix(h, 0)
	return x
}
