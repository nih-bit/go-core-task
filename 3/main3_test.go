package main

import (
	"maps"
	"testing"
)

func TestNewStringIntMap(t *testing.T) {
	m := NewStringIntMap()

	if m == nil {
		t.Fatal("NewStringIntMap() = nil, want указатель на структуру")
	}
	if got := m.Copy(); len(got) != 0 {
		t.Errorf("новая структура содержит %v, want пустую map", got)
	}
}

func TestStringIntMapAdd(t *testing.T) {
	tests := []struct {
		name   string
		keys   []string
		values []int
		want   map[string]int
	}{
		{"один элемент", []string{"a"}, []int{1}, map[string]int{"a": 1}},
		{"несколько элементов", []string{"a", "b", "c"}, []int{1, 2, 3}, map[string]int{"a": 1, "b": 2, "c": 3}},
		{"перезапись ключа", []string{"a", "a"}, []int{1, 2}, map[string]int{"a": 2}},
		{"нулевое значение", []string{"a"}, []int{0}, map[string]int{"a": 0}},
		{"отрицательное значение", []string{"a"}, []int{-5}, map[string]int{"a": -5}},
		{"пустой ключ", []string{""}, []int{1}, map[string]int{"": 1}},
		{"русский ключ", []string{"ключ"}, []int{1}, map[string]int{"ключ": 1}},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			m := NewStringIntMap()
			for i, key := range tc.keys {
				m.Add(key, tc.values[i])
			}

			if got := m.Copy(); !maps.Equal(got, tc.want) {
				t.Errorf("после Add структура содержит %v, want %v", got, tc.want)
			}
		})
	}
}

func TestStringIntMap_WithoutConstructor(t *testing.T) {
	var m StringIntMap

	m.Add("key", 42)

	if value, ok := m.Get("key"); !ok || value != 42 {
		t.Errorf("Get на нулевом значении = %d, %v; want 42, true", value, ok)
	}
}

func TestStringIntMapGet(t *testing.T) {
	m := NewStringIntMap()
	m.Add("zero", 0)
	m.Add("one", 1)

	tests := []struct {
		name      string
		key       string
		wantValue int
		wantOK    bool
	}{
		{"существующий ключ", "one", 1, true},
		{"существующий ключ с нулевым значением", "zero", 0, true},
		{"отсутствующий ключ", "missing", 0, false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			gotValue, gotOK := m.Get(tc.key)

			if gotValue != tc.wantValue || gotOK != tc.wantOK {
				t.Errorf("Get(%q) = %d, %v; want %d, %v", tc.key, gotValue, gotOK, tc.wantValue, tc.wantOK)
			}
		})
	}
}

func TestStringIntMapExists(t *testing.T) {
	m := NewStringIntMap()
	m.Add("a", 1)

	tests := []struct {
		name string
		key  string
		want bool
	}{
		{"существующий ключ", "a", true},
		{"отсутствующий ключ", "b", false},
		{"пустой ключ", "", false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			if got := m.Exists(tc.key); got != tc.want {
				t.Errorf("Exists(%q) = %v, want %v", tc.key, got, tc.want)
			}
		})
	}
}

func TestStringIntMapRemove(t *testing.T) {
	t.Run("существующий ключ", func(t *testing.T) {
		t.Parallel()

		m := NewStringIntMap()
		m.Add("a", 1)
		m.Add("b", 2)

		m.Remove("a")

		if got := m.Copy(); !maps.Equal(got, map[string]int{"b": 2}) {
			t.Errorf("после Remove структура содержит %v, want map[b:2]", got)
		}
	})

	t.Run("отсутствующий ключ", func(t *testing.T) {
		t.Parallel()

		m := NewStringIntMap()
		m.Add("a", 1)

		m.Remove("missing")

		if got := m.Copy(); !maps.Equal(got, map[string]int{"a": 1}) {
			t.Errorf("после Remove отсутствующего ключа структура содержит %v, want map[a:1]", got)
		}
	})

	t.Run("нулевое значение", func(t *testing.T) {
		t.Parallel()

		var m StringIntMap

		m.Remove("a")

		if got := m.Copy(); len(got) != 0 {
			t.Errorf("Copy после Remove на нулевом значении = %v, want пустую map", got)
		}
	})
}

func TestStringIntMapCopy(t *testing.T) {
	t.Run("содержит все элементы", func(t *testing.T) {
		t.Parallel()

		m := NewStringIntMap()
		m.Add("a", 1)
		m.Add("b", 2)

		if got := m.Copy(); !maps.Equal(got, map[string]int{"a": 1, "b": 2}) {
			t.Errorf("Copy() = %v, want map[a:1 b:2]", got)
		}
	})

	t.Run("изменение копии не влияет на оригинал", func(t *testing.T) {
		t.Parallel()

		m := NewStringIntMap()
		m.Add("a", 1)

		copied := m.Copy()
		copied["a"] = 999
		copied["b"] = 2

		if got := m.Copy(); !maps.Equal(got, map[string]int{"a": 1}) {
			t.Errorf("после правки копии структура содержит %v, want map[a:1]", got)
		}
	})

	t.Run("изменение оригинала не влияет на копию", func(t *testing.T) {
		t.Parallel()

		m := NewStringIntMap()
		m.Add("a", 1)

		copied := m.Copy()
		m.Add("a", 999)

		if !maps.Equal(copied, map[string]int{"a": 1}) {
			t.Errorf("после правки оригинала копия содержит %v, want map[a:1]", copied)
		}
	})

	t.Run("нулевое значение даёт не-nil пустую map", func(t *testing.T) {
		t.Parallel()

		var m StringIntMap

		got := m.Copy()

		if got == nil {
			t.Fatal("Copy() = nil, want не-nil пустую map")
		}
		if len(got) != 0 {
			t.Errorf("Copy() = %v, want пустую map", got)
		}
	})
}
