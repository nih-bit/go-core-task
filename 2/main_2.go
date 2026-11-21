package main

import (
	"fmt"
	"math/rand/v2"
	"slices"
)

func main() {
	const size = 10
	originalSlice := make([]int, 0, size)

	for i := 0; i < size; i++ {
		originalSlice = append(originalSlice, rand.IntN(1000))
	}

	fmt.Printf("original slice: %v\n", originalSlice)

	fmt.Printf("even slice: %v\n", sliceExample(originalSlice))
	fmt.Printf("original slice (unchanged): %v\n", originalSlice)

	fmt.Printf("slice with new element: %v\n", addElements(originalSlice, rand.IntN(1000)))
	fmt.Printf("original slice (unchanged): %v\n", originalSlice)

	newSlice := copySlice(originalSlice)
	newSlice[0] = -1
	newSlice[len(newSlice)-1] = -2
	fmt.Printf("copy after copy[0] = -1, copy[last] = -2: %v\n", newSlice)
	fmt.Printf("original slice (unchanged): %v\n", originalSlice)

	const indexToRemove = 5
	fmt.Printf("slice without element by index %d: %v\n", indexToRemove, removeElement(originalSlice, indexToRemove))
	fmt.Printf("original slice (unchanged): %v\n", originalSlice)
}

func sliceExample(slice []int) []int {
	resultSlice := make([]int, 0, len(slice))

	for _, n := range slice {
		if n%2 == 0 {
			resultSlice = append(resultSlice, n)
		}
	}

	return resultSlice
}

func copySlice(original []int) []int {
	return slices.Clone(original)
}

func addElements(slice []int, n int) []int {
	return append(copySlice(slice), n)
}

func removeElement(slice []int, index int) []int {
	if index < 0 || index >= len(slice) {
		return copySlice(slice)
	}

	resultSlice := make([]int, 0, len(slice)-1)

	for i, n := range slice {
		if i != index {
			resultSlice = append(resultSlice, n)
		}
	}

	return resultSlice
}
