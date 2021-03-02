// Copyright 2021 lkevinzc. All rights reserved.

// Package heap provides basic heap data structure to facilitates
// priority queue implementation. A heap is a tree with the property
// that each node is the minimum-valued node in its subtree.
//
// The minimum element in the tree is the root, at index **1**, which
// makes the indexing a bit easier.
//
// This implementation provides the option to record the item creation
// time, so that the Less() compares the time if there is a tie in the
// priority. This is useful for dealing with requests (FCFS).
//
package heap

import "time"

// An Item contains any data with a priority value.
type Item struct {
	Priority  int
	Data      interface{}
	CreatedAt time.Time
}

// ItemHeap implements the basic min heap of Item.
type ItemHeap []*Item

// NewHeap returns a ItemHeap instance that has a dummy first item for
// easier indexing.
func NewHeap() ItemHeap {
	h := ItemHeap{&Item{
		Priority: 0,
		Data:     nil,
	}}
	return h
}

// Len returns heap size (n-1) instead of the real array size (n).
func (h ItemHeap) Len() int {
	return len(h) - 1
}

// Empty tests if the heap (not underlying array) is empty.
func (h ItemHeap) Empty() bool {
	return h.Len() == 0
}

// Less serves as a comparator.
func (h ItemHeap) Less(i, j int) bool {
	if h[i].Priority == h[j].Priority {
		return h[i].CreatedAt.Before(h[j].CreatedAt)
	}
	return h[i].Priority < h[j].Priority
}

// Swap swaps two array elements (i.e. items).
func (h ItemHeap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}

// Push pushes the element x onto the heap.
// The complexity is O(log n) where n = h.Len().
func (h *ItemHeap) Push(x interface{}) {
	item := x.(*Item)
	*h = append(*h, item)
	h.up(h.Len())
}

// Pop removes and returns the minimum element (according to Less) from the heap.
// The complexity is O(log n) where n = h.Len().
// If the heap is empty, Pop returns nil.
func (h *ItemHeap) Pop() interface{} {
	if h.Empty() {
		return nil
	}
	n := h.Len()
	h.Swap(1, n) // item at index 1 is the valid smallest
	old := *h
	item := old[n]
	old[n] = nil // avoid memory leak
	*h = old[0:n]
	h.down(1)
	return item
}

func (h *ItemHeap) up(j int) {
	i := parent(j)
	if j > 1 && h.Less(j, i) {
		h.Swap(i, j)
		h.up(i)
	}
}

func (h *ItemHeap) down(j int) {
	n := h.Len()
	l := leftChild(j)
	r := rightChild(j)
	if l > n {
		return
	}
	var smallestChild int
	if r > n {
		smallestChild = l
	} else {
		if h.Less(l, r) {
			smallestChild = l
		} else {
			smallestChild = r
		}
	}
	if smallestChild <= n && h.Less(smallestChild, j) {
		h.Swap(j, smallestChild)
		h.down(smallestChild)
	}
}

func parent(k int) int     { return k / 2 }
func leftChild(k int) int  { return k * 2 }
func rightChild(k int) int { return k*2 + 1 }
