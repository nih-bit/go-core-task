package main

import "fmt"

type StringIntMap struct {
	innerMap map[string]int
}

func NewStringIntMap() *StringIntMap {
	return &StringIntMap{innerMap: make(map[string]int)}
}

func (sim *StringIntMap) Add(key string, value int) {
	if sim.innerMap == nil {
		sim.innerMap = make(map[string]int)
	}

	sim.innerMap[key] = value
}

func (sim *StringIntMap) Remove(key string) {
	delete(sim.innerMap, key)
}

func (sim *StringIntMap) Copy() map[string]int {
	newMap := make(map[string]int, len(sim.innerMap))

	for key, value := range sim.innerMap {
		newMap[key] = value
	}

	return newMap
}

func (sim *StringIntMap) Exists(key string) bool {
	_, ok := sim.innerMap[key]

	return ok
}

func (sim *StringIntMap) Get(key string) (int, bool) {
	value, ok := sim.innerMap[key]

	return value, ok
}

func main() {
	m := NewStringIntMap()

	m.Add("one", 1)
	m.Add("two", 2)
	m.Add("three", 3)
	fmt.Printf("after Add: %v\n", m.Copy())

	m.Remove("two")
	fmt.Printf("after Remove(\"two\"): %v\n", m.Copy())

	fmt.Printf("Exists(\"one\"): %v\n", m.Exists("one"))
	fmt.Printf("Exists(\"two\"): %v\n", m.Exists("two"))

	value, ok := m.Get("three")
	fmt.Printf("Get(\"three\"): %d, %v\n", value, ok)

	value, ok = m.Get("two")
	fmt.Printf("Get(\"two\"): %d, %v\n", value, ok)
}
