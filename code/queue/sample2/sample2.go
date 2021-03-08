package main

import (
	"log"
	"runtime"
	"sync"
	"time"
)

// Queue: FIFO
type Queue struct {
	q chan int
}

// New: create a new instance
func New(size int) *Queue {
	return &Queue{make(chan int, size)}
}

// Add: enqueue
func (q *Queue) Add(s int) {
	q.q <- s
}

// Get: dequeue
func (q *Queue) Get() (int, bool) {
	n, ok := <-q.q
	return n, ok
}

//Complete: close channel
func (q *Queue) Complete() {
	close(q.q)
}

const MaxWorkers = 2

func Workers(wg *sync.WaitGroup, q *Queue) {
	for i := 0; i < MaxWorkers; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			for {
				if n, ok := q.Get(); ok {
					log.Printf("Worker(%d): start Task(%d)\n", i, n)
					time.Sleep(2 * time.Second) //working...
					log.Printf("Worker(%d): end Task(%d)\n", i, n)
				} else {
					break
				}
			}
			log.Printf("Worker(%d): return home\n", i)
		}(i + 1)
	}
}

func Manager(tasklist []int) *Queue {
	plan := New(MaxWorkers)
	go func() {
		defer plan.Complete()
		for _, n := range tasklist {
			plan.Add(n)
			log.Printf("Manager: set Task(%d)\n", n)
		}
		log.Println("Manager: return home")
	}()
	return plan
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	tasklist := []int{1, 2, 3, 4, 5}
	log.Println("Start...")
	var wg sync.WaitGroup
	plan := Manager(tasklist)
	Workers(&wg, plan)
	wg.Wait()
	log.Println("...End")
}
