package main

import "fmt"

func main() {
	a := []int{65, 3, 58, 678, 64}
	b := []int{64, 2, 3, 43}
	ok, result := intersections(a, b)

	fmt.Printf("first slice: %v\nsecond slice: %v\nresult: %t, %v\n", a, b, ok, result)
}

func intersections(slice1, slice2 []int) (bool, []int) {
	seen := make(map[int]struct{}, len(slice1))

	for _, n := range slice1 {
		seen[n] = struct{}{}
	}

	result := make([]int, 0, min(len(slice1), len(slice2)))

	for _, n := range slice2 {
		if _, ok := seen[n]; ok {
			result = append(result, n)
			delete(seen, n)
		}
	}

	return len(result) > 0, result
}
