package main

import (
	"slices"
	"testing"
)

func TestSliceExample(t *testing.T) {
	tests := []struct {
		name  string
		input []int
		want  []int
	}{
		{"смешанные числа", []int{1, 2, 3, 4, 5, 6}, []int{2, 4, 6}},
		{"только чётные", []int{2, 4, 6}, []int{2, 4, 6}},
		{"только нечётные", []int{1, 3, 5}, []int{}},
		{"порядок сохраняется", []int{8, 1, 6, 3, 4}, []int{8, 6, 4}},
		{"отрицательные числа", []int{-4, -3, -2, -1, 0}, []int{-4, -2, 0}},
		{"пустой срез", []int{}, []int{}},
		{"nil", nil, []int{}},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			before := slices.Clone(tc.input)

			got := sliceExample(tc.input)

			if !slices.Equal(got, tc.want) {
				t.Errorf("sliceExample(%v) = %v, want %v", tc.input, got, tc.want)
			}
			if !slices.Equal(tc.input, before) {
				t.Errorf("исходный срез изменился: было %v, стало %v", before, tc.input)
			}
		})
	}
}

func TestSliceExample_ReturnsNewSlice(t *testing.T) {
	input := []int{2, 4, 6}

	got := sliceExample(input)
	got[0] = -1

	if !slices.Equal(input, []int{2, 4, 6}) {
		t.Errorf("изменение результата изменило вход: input = %v, want [2 4 6]", input)
	}
}

func TestAddElements(t *testing.T) {
	tests := []struct {
		name  string
		input []int
		n     int
		want  []int
	}{
		{"непустой срез", []int{1, 2, 3}, 4, []int{1, 2, 3, 4}},
		{"пустой срез", []int{}, 1, []int{1}},
		{"nil", nil, 7, []int{7}},
		{"отрицательный элемент", []int{1}, -5, []int{1, -5}},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			before := slices.Clone(tc.input)

			got := addElements(tc.input, tc.n)

			if !slices.Equal(got, tc.want) {
				t.Errorf("addElements(%v, %d) = %v, want %v", tc.input, tc.n, got, tc.want)
			}
			if !slices.Equal(tc.input, before) {
				t.Errorf("исходный срез изменился: было %v, стало %v", before, tc.input)
			}
		})
	}
}

func TestAddElements_ReturnsNewSlice(t *testing.T) {
	input := []int{1, 2, 3}

	got := addElements(input, 4)

	got[0] = -1

	if !slices.Equal(input, []int{1, 2, 3}) {
		t.Errorf("изменение результата изменило вход: input = %v, want [1 2 3]", input)
	}
}

func TestCopySlice(t *testing.T) {
	t.Run("содержимое совпадает", func(t *testing.T) {
		t.Parallel()

		original := []int{1, 2, 3, 4, 5}

		if got := copySlice(original); !slices.Equal(got, original) {
			t.Errorf("copySlice(%v) = %v, want %v", original, got, original)
		}
	})

	t.Run("изменение копии не влияет на оригинал", func(t *testing.T) {
		t.Parallel()

		original := []int{1, 2, 3, 4, 5}

		got := copySlice(original)
		got[0] = -1

		if original[0] != 1 {
			t.Errorf("правка копии изменила оригинал: original[0] = %d, want 1", original[0])
		}
	})

	t.Run("изменение оригинала не влияет на копию", func(t *testing.T) {
		t.Parallel()

		original := []int{1, 2, 3, 4, 5}

		got := copySlice(original)
		original[0] = -1

		if got[0] != 1 {
			t.Errorf("правка оригинала изменила копию: copy[0] = %d, want 1", got[0])
		}
	})

	t.Run("nil", func(t *testing.T) {
		t.Parallel()

		if got := copySlice(nil); len(got) != 0 {
			t.Errorf("copySlice(nil) = %v, want пустой срез", got)
		}
	})

	t.Run("пустой срез", func(t *testing.T) {
		t.Parallel()

		if got := copySlice([]int{}); len(got) != 0 {
			t.Errorf("copySlice([]) = %v, want пустой срез", got)
		}
	})
}

func TestRemoveElement(t *testing.T) {
	tests := []struct {
		name  string
		input []int
		index int
		want  []int
	}{
		{"первый элемент", []int{1, 2, 3}, 0, []int{2, 3}},
		{"середина", []int{1, 2, 3}, 1, []int{1, 3}},
		{"последний элемент", []int{1, 2, 3}, 2, []int{1, 2}},
		{"единственный элемент", []int{7}, 0, []int{}},
		{"индекс == len", []int{1, 2, 3}, 3, []int{1, 2, 3}},
		{"индекс > len", []int{1, 2, 3}, 99, []int{1, 2, 3}},
		{"отрицательный индекс", []int{1, 2, 3}, -1, []int{1, 2, 3}},
		{"пустой срез", []int{}, 0, []int{}},
		{"пустой срез, отрицательный индекс", []int{}, -1, []int{}},
		{"nil", nil, 0, nil},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			before := slices.Clone(tc.input)

			got := removeElement(tc.input, tc.index)

			if !slices.Equal(got, tc.want) {
				t.Errorf("removeElement(%v, %d) = %v, want %v", tc.input, tc.index, got, tc.want)
			}
			if !slices.Equal(tc.input, before) {
				t.Errorf("исходный срез изменился: было %v, стало %v", before, tc.input)
			}
		})
	}
}

func TestRemoveElement_ReturnsNewSlice(t *testing.T) {
	input := []int{1, 2, 3}

	got := removeElement(input, 0)
	got[0] = -1

	if !slices.Equal(input, []int{1, 2, 3}) {
		t.Errorf("изменение результата изменило вход: input = %v, want [1 2 3]", input)
	}
}
