// Copyright 2021 lkevinzc. All rights reserved.

package requestpq

import (
	"fmt"
	"math/rand"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func verify(t *testing.T, q *Queue) {
	if q.Empty() {
		return
	}
	prevData, err := q.Dequeue()
	assert.Equal(t, nil, err)
	for !q.Empty() {
		data, err := q.Dequeue()
		assert.Equal(t, nil, err)
		assert.LessOrEqual(t, prevData, data)
		prevData = data
	}
}

func isAscending(t *testing.T, arr []interface{}) {
	i := 0
	for j := 1; j < len(arr); j++ {
		assert.LessOrEqual(t, arr[i], arr[j])
		i = j
	}
}

func TestNewQueue(t *testing.T) {
	q := NewQueue()
	assert.Equal(t, 0, q.Len())
	q.Enqueue(`test`, 20)
	assert.Equal(t, 1, q.Len())
	q.Enqueue(3.14, 19)
	assert.Equal(t, 2, q.Len())
	for !q.Empty() {
		data, _ := q.Dequeue()
		fmt.Println(data) // heterogeneous data types
	}
}

func TestQueue(t *testing.T) {
	t.Run("random priority, more enqueue than dequeue", func(t *testing.T) {
		q := NewQueue()
		for i := 0; i < 10000; i++ {
			if rand.Intn(3) != 0 { // enq with prob = 2/3
				v := rand.Intn(20)
				q.Enqueue(v, v)
			} else {
				_, err := q.Dequeue()
				if err != nil {
					assert.Equal(t, true, q.Empty())
				}
			}
		}
		verify(t, q)
	})

	t.Run("random priority, more dequeue than enqueue", func(t *testing.T) {
		q := NewQueue()
		for i := 0; i < 10000; i++ {
			if rand.Intn(3) == 0 { // enq with prob = 1/3
				v := rand.Intn(20)
				q.Enqueue(v, v)
			} else {
				_, err := q.Dequeue()
				if err != nil {
					assert.Equal(t, true, q.Empty())
				}
			}
		}
		verify(t, q)
	})

	t.Run("equal priority, test enqueue sequence", func(t *testing.T) {
		q := NewQueue()
		for i := 0; i < 10000; i++ {
			if rand.Intn(3) == 0 { // enq with prob = 1/3
				q.Enqueue(time.Now().UnixNano(), 20)
			} else {
				_, err := q.Dequeue()
				if err != nil {
					assert.Equal(t, true, q.Empty())
				}
			}
		}
		verify(t, q)
	})
}

func TestDecorateChannel(t *testing.T) {
	t.Run("enqueue-dequeue test", func(t *testing.T) {
		N := 100
		inChan := make(chan *Task, 1)
		outChan := DecorateChannel(inChan, 1)
		var wg sync.WaitGroup
		wg.Add(1)
		i := 0
		go func() {
			for {
				<-outChan
				i++
				if i == N {
					wg.Done()
					break
				}
			}
		}()
		for i := 0; i < N; i++ {
			inChan <- &Task{
				Data:     i,
				Priority: i,
			}
		}
	})

	t.Run("random priority for sanity check", func(t *testing.T) {
		N := 5000
		buffer := N
		inChan := make(chan *Task, buffer)
		outChan := DecorateChannel(inChan, buffer)
		blocker := make(chan bool)
		i := 0
		go func() { // producer
			for i := 0; i < N; i++ {
				v := rand.Intn(20)
				inChan <- &Task{
					Data:     v,
					Priority: v,
				}
			}
			blocker <- true
		}()
		<-blocker
		var localArr []interface{}
		for {
			data := <-outChan
			localArr = append(localArr, data)
			i++
			if i == N {
				break
			}
		}
		for i := 0; i < 20; i++ {
			fmt.Printf("%v ", localArr[i])
		}
		//isAscending(t, localArr)
	})
}

// go test -v -race -cover
// go test -bench=.
