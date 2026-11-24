package main

import (
	"context"
	"fmt"
)

func producer(ctx context.Context, numbers ...uint8) <-chan uint8 {
	ch := make(chan uint8)

	go func() {
		defer close(ch)

		for _, n := range numbers {
			select {
			case ch <- n:
			case <-ctx.Done():
				return
			}
		}
	}()

	return ch
}

func cubeFloat64(ctx context.Context, in <-chan uint8) <-chan float64 {
	out := make(chan float64)

	go func() {
		defer close(out)

		for {
			select {
			case n, ok := <-in:
				if !ok {
					return
				}

				f := float64(n)
				cube := f * f * f

				select {
				case out <- cube:
				case <-ctx.Done():
					return
				}
			case <-ctx.Done():
				return
			}
		}
	}()

	return out
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	for f := range cubeFloat64(ctx, producer(ctx, 0, 1, 3, 7, 15, 31, 63, 127, 255)) {
		fmt.Printf("receive %.2f from pipeline\n", f)
	}
}
