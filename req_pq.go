package requestpq

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/lkevinzc/requestpq/heap"
)

// Task defines the input format of decorated channel.
type Task struct {
	Data     interface{}
	Priority int
}

// Queue is a thread-safe priority queue.
type Queue struct {
	heap *heap.ItemHeap
	lock sync.Mutex
}

// NewQueue is the constructor of Queue.
func NewQueue() *Queue {
	h := heap.NewHeap()
	q := Queue{heap: &h}
	return &q
}

// Enqueue puts the data into the priority queue with a timestamp.
func (q *Queue) Enqueue(data interface{}, priority int) {
	q.lock.Lock()
	defer q.lock.Unlock()
	item := heap.Item{
		Priority:  priority,
		Data:      data,
		CreatedAt: time.Now(),
	}
	q.heap.Push(&item)
}

// Dequeue gets & removes the data with highest priority from the queue.
func (q *Queue) Dequeue() (interface{}, error) {
	q.lock.Lock()
	defer q.lock.Unlock()
	item := q.heap.Pop()
	if item == nil {
		return nil, errors.New("pop an empty queue")
	}
	return item.(*heap.Item).Data, nil
}

// Len returns the size of the priority queue.
func (q *Queue) Len() int {
	q.lock.Lock()
	defer q.lock.Unlock()
	return q.heap.Len()
}

// Empty tests if the queue is empty.
func (q *Queue) Empty() bool {
	q.lock.Lock()
	defer q.lock.Unlock()
	return q.heap.Empty()
}

// DecorateChannel transforms a FIFO queue of normal channel
// into priority queue with decorated channel.
func DecorateChannel(inChan chan *Task) (outChan chan interface{}) {
	outChan = make(chan interface{})
	pq := NewQueue()
	cond := sync.NewCond(&pq.lock)
	go func() {
		for task := range inChan {
			pq.Enqueue(task.Data, task.Priority)
			cond.Signal()
		}
	}()
	go func() {
		for {
			pq.lock.Lock()
			if pq.heap.Empty() {
				cond.Wait()
			}
			item := pq.heap.Pop()
			if item == nil {
				panic(fmt.Sprintf("pop an empty queue"))
			}
			data := item.(*heap.Item).Data
			pq.lock.Unlock()
			outChan <- data
		}
	}()
	return
}
