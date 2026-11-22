package main

import "fmt"

func main() {
	slice1 := []string{"apple", "banana", "cherry", "date", "43", "lead", "gno1"}
	slice2 := []string{"banana", "date", "fig"}
	result := diff(slice1, slice2)

	fmt.Printf("first slice: %v\nsecond slice: %v\nresult: %v\n", slice1, slice2, result)
}

func diff(slice1, slice2 []string) []string {
	seen := make(map[string]struct{}, len(slice2))

	for _, s := range slice2 {
		seen[s] = struct{}{}
	}

	result := make([]string, 0, len(slice1))

	for _, s := range slice1 {
		if _, ok := seen[s]; !ok {
			result = append(result, s)
		}
	}

	return result
}
