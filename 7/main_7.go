package main

import (
	"fmt"
	"math/rand/v2"
	"sync"
)

func main() {
	for n := range merge(randInts(2, 10), randInts(3, 100), randInts(4, 1000)) {
		fmt.Println(n)
	}
}

func merge(chs ...<-chan int) <-chan int {
	out := make(chan int)
	var wg sync.WaitGroup

	for _, ch := range chs {
		if ch == nil {
			continue
		}

		wg.Add(1)

		go func(ch <-chan int) {
			defer wg.Done()

			for value := range ch {
				out <- value
			}
		}(ch)
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}

func randInts(n int, maxInt int) <-chan int {
	ch := make(chan int)

	go func() {
		defer close(ch)

		for range n {
			ch <- rand.IntN(maxInt)
		}
	}()

	return ch
}
