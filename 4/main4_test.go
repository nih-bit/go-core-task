package main

import (
	"slices"
	"testing"
)

func TestDiff(t *testing.T) {
	tests := []struct {
		name   string
		first  []string
		second []string
		want   []string
	}{
		{
			name:   "пример из задания",
			first:  []string{"apple", "banana", "cherry", "date", "43", "lead", "gno1"},
			second: []string{"banana", "date", "fig"},
			want:   []string{"apple", "cherry", "43", "lead", "gno1"},
		},
		{
			name:   "нет пересечения",
			first:  []string{"a", "b"},
			second: []string{"c", "d"},
			want:   []string{"a", "b"},
		},
		{
			name:   "полное пересечение",
			first:  []string{"a", "b"},
			second: []string{"a", "b"},
			want:   []string{},
		},
		{
			name:   "второй слайс — надмножество",
			first:  []string{"a"},
			second: []string{"a", "b", "c"},
			want:   []string{},
		},
		{
			name:   "оба пустые",
			first:  nil,
			second: nil,
			want:   []string{},
		},
		{
			name:   "первый пустой",
			first:  nil,
			second: []string{"a"},
			want:   []string{},
		},
		{
			name:   "второй пустой",
			first:  []string{"a", "b"},
			second: nil,
			want:   []string{"a", "b"},
		},
		{
			name:   "дубликаты в первом слайсе сохраняются",
			first:  []string{"a", "a", "b"},
			second: []string{"b"},
			want:   []string{"a", "a"},
		},
		{
			name:   "дубликаты во втором слайсе не важны",
			first:  []string{"a", "b"},
			second: []string{"b", "b", "b"},
			want:   []string{"a"},
		},
		{
			name:   "регистр значим",
			first:  []string{"Apple", "apple"},
			second: []string{"apple"},
			want:   []string{"Apple"},
		},
		{
			name:   "пустая строка",
			first:  []string{"", "a"},
			second: []string{"a"},
			want:   []string{""},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got := diff(tc.first, tc.second)

			if !slices.Equal(got, tc.want) {
				t.Errorf("diff(%v, %v) = %v, want %v", tc.first, tc.second, got, tc.want)
			}
			if got == nil {
				t.Error("diff вернул nil, want не-nil срез")
			}
		})
	}
}
