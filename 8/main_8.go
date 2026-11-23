package main

import (
	"fmt"
	"sync"
	"time"
)

type WaitGroup struct {
	mu      sync.Mutex
	counter int
	done    chan struct{}
}

func (wg *WaitGroup) Add(delta int) {
	wg.mu.Lock()
	defer wg.mu.Unlock()

	if delta < 0 {
		panic("negative delta")
	}

	if wg.counter == 0 && delta > 0 {
		wg.done = make(chan struct{})
	}

	wg.counter += delta
}

func (wg *WaitGroup) Done() {
	wg.mu.Lock()
	defer wg.mu.Unlock()

	if wg.counter == 0 {
		panic("Done without Add")
	}

	wg.counter--

	if wg.counter == 0 {
		close(wg.done)
	}
}

func (wg *WaitGroup) Wait() {
	wg.mu.Lock()

	ch := wg.done

	wg.mu.Unlock()

	if ch != nil {
		<-ch
	}
}

func main() {
	var wg WaitGroup

	for i := 0; i < 10; i++ {
		wg.Add(1)

		go func(n int) {
			defer wg.Done()

			fmt.Printf("goroutine number %d\n", n)

			<-time.After(time.Second)
		}(i + 1)
	}

	wg.Wait()

	fmt.Println("all goroutine finished after wait")
}
