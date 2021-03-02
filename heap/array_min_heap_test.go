// Copyright 2021 lkevinzc. All rights reserved.
// Adapted from go src files (src/container/heap_test.go).

package heap

import (
	"fmt"
	"math/rand"
	"sync"
	"testing"
	"time"
)

func (h ItemHeap) verify(t *testing.T, i int) {
	t.Helper()
	n := h.Len()
	j1 := 2 * i
	j2 := 2*i + 1
	if j1 < n {
		if h.Less(j1, i) {
			t.Errorf("heap invariant invalidated [%d] = %v > [%d] = %v", i, h[i], j1, h[j1])
			return
		}
		h.verify(t, j1)
	}
	if j2 < n {
		if h.Less(j2, i) {
			t.Errorf("heap invariant invalidated [%d] = %v > [%d] = %v", i, h[i], j1, h[j2])
			return
		}
		h.verify(t, j2)
	}
}

func TestInit0(t *testing.T) {
	h := NewHeap()
	for i := 20; i > 0; i-- {
		h.Push(&Item{
			Priority:  0,
			Data:      `test`,
			CreatedAt: time.Now(),
		}) // all elements are the same
	}

	h.verify(t, 1)

	for i := 1; h.Len() > 0; i++ {
		x := h.Pop().(*Item)
		if x.Priority != 0 {
			t.Errorf("%d.th pop got %v; want %d", i, x, 0)
		}
	}
}

func TestInit1(t *testing.T) {
	h := NewHeap()
	for i := 20; i > 0; i-- {
		h.Push(&Item{
			Priority:  i,
			Data:      `test`,
			CreatedAt: time.Now(),
		}) // all elements are different
	}

	h.verify(t, 1)

	for i := 1; h.Len() > 0; i++ {
		x := h.Pop().(*Item)
		h.verify(t, 1)
		if x.Priority != i {
			t.Errorf("%d.th pop got %v; want %d", i, x, i)
		}
	}
}

func TestOrder(t *testing.T) {
	h := NewHeap()
	h.verify(t, 1)

	for i := 20; i > 10; i-- {
		h.Push(&Item{
			Priority:  i,
			Data:      `test`,
			CreatedAt: time.Now(),
		})
	}

	h.verify(t, 1)

	for i := 10; i > 0; i-- {
		h.Push(&Item{
			Priority:  i,
			Data:      `test`,
			CreatedAt: time.Now(),
		})
		h.verify(t, 1)
	}

	for i := 1; h.Len() > 0; i++ {
		x := h.Pop().(*Item)
		if i < 20 {
			h.Push(&Item{
				Priority:  20 + i,
				Data:      `test`,
				CreatedAt: time.Now(),
			})
		}
		h.verify(t, 1)
		if x.Priority != i {
			t.Errorf("%d.th pop got %v; want %d", i, x, i)
		}
	}
}

func TestRandom(t *testing.T) {
	h := NewHeap()
	h.verify(t, 1)

	for i := 0; i < 100; i++ {
		h.Push(&Item{
			Priority:  rand.Intn(20),
			Data:      `test`,
			CreatedAt: time.Now(),
		})
	}

	h.verify(t, 1)

	for j := 10; j > 0; j-- {
		_ = h.Pop().(*Item)
		h.verify(t, 1)
	}
}

func TestRandomVisualize(t *testing.T) {
	h := NewHeap()
	h.verify(t, 1)

	for i := 0; i < 20; i++ {
		h.Push(&Item{
			Priority:  rand.Intn(20),
			Data:      `test`,
			CreatedAt: time.Now(),
		})
	}

	h.verify(t, 1)

	for !h.Empty() {
		x := h.Pop().(*Item)
		fmt.Printf("%v ", x.Priority)
	}
	fmt.Println()
}

func TestEqualPriorityNoTime(t *testing.T) {
	h := NewHeap()
	h.verify(t, 1)

	for i := 0; i < 20; i++ {
		h.Push(&Item{
			Priority: 20,
			Data:     fmt.Sprintf("test%v", i),
		})
	}

	h.verify(t, 1)
	t.Logf("The following sequence is out of order.")
	for !h.Empty() {
		x := h.Pop().(*Item)
		fmt.Printf("<%v %v>", x.Priority, x.Data)
	}
	fmt.Println()
}

func TestEqualPriority(t *testing.T) {
	h := NewHeap()
	h.verify(t, 1)

	for i := 0; i < 20; i++ {
		h.Push(&Item{
			Priority:  20,
			Data:      fmt.Sprintf("test%v", i),
			CreatedAt: time.Now(),
		})
	}

	h.verify(t, 1)
	t.Logf("The following sequence is in the order of insertion.")
	for !h.Empty() {
		x := h.Pop().(*Item)
		fmt.Printf("<%v %v>", x.Priority, x.Data)
	}
	fmt.Println()
}

func TestPopEmpty(t *testing.T) {
	h := NewHeap()
	h.verify(t, 1)
	for i := 0; i < 5; i++ {
		h.Push(&Item{
			Priority:  i,
			Data:      `test`,
			CreatedAt: time.Now(),
		})
	}

	for i := 0; i < 7; i++ {
		y := h.Pop()
		if y != nil {
			x := y.(*Item)
			if x.Priority != i {
				t.Errorf("%d.th pop got %v; want %d", i, x, i)
			}
		}
		h.verify(t, 1)
	}
}

func BenchmarkHeapDup(b *testing.B) {
	const n = 10000
	h := NewHeap()
	for i := 0; i < b.N; i++ {
		for j := 0; j < n; j++ {
			h.Push(&Item{
				Priority:  0,
				Data:      `test`,
				CreatedAt: time.Now(),
			}) // all elements are the same
		}
		for h.Len() > 0 {
			h.Pop()
		}
	}
}

func BenchmarkHeapDupNoTime(b *testing.B) {
	const n = 10000
	h := NewHeap()
	for i := 0; i < b.N; i++ {
		for j := 0; j < n; j++ {
			h.Push(&Item{
				Priority: 0,
				Data:     `test`,
			}) // all elements are the same
		}
		for h.Len() > 0 {
			h.Pop()
		}
	}
}

func BenchmarkHeapRnd(b *testing.B) {
	const n = 10000
	h := NewHeap()
	for i := 0; i < b.N; i++ {
		for j := 0; j < n; j++ {
			h.Push(&Item{
				Priority:  rand.Intn(20),
				Data:      `test`,
				CreatedAt: time.Now(),
			}) // all elements are random
		}
		for h.Len() > 0 {
			h.Pop()
		}
	}
}

func BenchmarkHeapRndNoTime(b *testing.B) {
	const n = 10000
	h := NewHeap()
	for i := 0; i < b.N; i++ {
		for j := 0; j < n; j++ {
			h.Push(&Item{
				Priority: rand.Intn(20),
				Data:     `test`,
			}) // all elements are random
		}
		for h.Len() > 0 {
			h.Pop()
		}
	}
}

func BenchmarkChanQDup(b *testing.B) {
	const n = 10000
	ch := make(chan interface{}, b.N)
	var wg sync.WaitGroup
	wg.Add(1)
	i := 0
	go func() {
		for {
			for i := 0; i < n; i++ {
				<-ch
			}
			i++
			if i == b.N {
				wg.Done()
				break
			}
		}
	}()

	for i := 0; i < b.N; i++ {
		for j := 0; j < n; j++ {
			ch <- &Item{
				Priority:  rand.Intn(20),
				Data:      `test`,
				CreatedAt: time.Now(),
			}
		}
	}
}

func BenchmarkChanQDupNoTime(b *testing.B) {
	const n = 10000
	ch := make(chan interface{}, b.N)
	var wg sync.WaitGroup
	wg.Add(1)
	i := 0
	go func() {
		for {
			for i := 0; i < n; i++ {
				<-ch
			}
			i++
			if i == b.N {
				wg.Done()
				break
			}
		}
	}()

	for i := 0; i < b.N; i++ {
		for j := 0; j < n; j++ {
			ch <- &Item{
				Priority: rand.Intn(20),
				Data:     `test`,
			}
		}
	}
}

// go test -v -race -cover
// go test -bench=.
