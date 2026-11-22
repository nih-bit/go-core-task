package main

import (
	"fmt"
	"math/rand/v2"
)

func main() {
	for n := range randInts(10, 1000) {
		fmt.Println(n)
	}
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
