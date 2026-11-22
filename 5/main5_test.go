package main

import (
	"slices"
	"testing"
)

func TestIntersections(t *testing.T) {
	tests := []struct {
		name       string
		first      []int
		second     []int
		wantOK     bool
		wantResult []int
	}{
		{
			name:       "пример из задания",
			first:      []int{65, 3, 58, 678, 64},
			second:     []int{64, 2, 3, 43},
			wantOK:     true,
			wantResult: []int{64, 3},
		},
		{
			name:       "нет пересечения",
			first:      []int{65, 58, 678},
			second:     []int{64, 3, 43},
			wantOK:     false,
			wantResult: []int{},
		},
		{
			name:       "полное пересечение",
			first:      []int{2, 3, 43, 64},
			second:     []int{64, 2, 3, 43},
			wantOK:     true,
			wantResult: []int{64, 2, 3, 43},
		},
		{
			name:       "оба пустые",
			first:      nil,
			second:     nil,
			wantOK:     false,
			wantResult: []int{},
		},
		{
			name:       "первый пустой",
			first:      nil,
			second:     []int{64, 2, 3, 43},
			wantOK:     false,
			wantResult: []int{},
		},
		{
			name:       "второй пустой",
			first:      []int{65, 3, 58, 678, 64},
			second:     nil,
			wantOK:     false,
			wantResult: []int{},
		},
		{
			name:       "дубликаты в первом слайсе не сохраняются",
			first:      []int{1, 1, 2, 2, 3, 3, 4, 4},
			second:     []int{1, 3},
			wantOK:     true,
			wantResult: []int{1, 3},
		},
		{
			name:       "дубликаты во втором слайсе не сохраняются",
			first:      []int{1, 3},
			second:     []int{1, 1, 2, 2, 3, 3, 4, 4},
			wantOK:     true,
			wantResult: []int{1, 3},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			gotOK, gotResult := intersections(tc.first, tc.second)

			if gotOK != tc.wantOK || !slices.Equal(gotResult, tc.wantResult) {
				t.Errorf("intersections(%v, %v) = %t, %v, want %t, %v", tc.first, tc.second, gotOK, gotResult, tc.wantOK, tc.wantResult)
			}

			if gotResult == nil {
				t.Error("intersections вернул nil, want не-nil срез")
			}
		})
	}
}
