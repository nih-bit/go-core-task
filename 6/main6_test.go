package main

import (
	"fmt"
	"slices"
	"testing"
	"time"
)

// collect вычитывает канал до закрытия и возвращает все значения
func collect(ch <-chan int, want int) ([]int, error) {
	timeout := time.After(2 * time.Second)
	got := make([]int, 0, want)

	for {
		select {
		case value, ok := <-ch:
			if !ok {
				return got, nil
			}

			got = append(got, value)

			if len(got) > want {
				return got, fmt.Errorf("получено больше %d значений — канал не закрывается", want)
			}
		case <-timeout:
			return got, fmt.Errorf("таймаут: получено %d значений из %d", len(got), want)
		}
	}
}

func TestRandInts(t *testing.T) {
	tests := []struct {
		name   string
		n      int
		maxInt int
	}{
		{"десять чисел из [0, 1000)", 10, 1000},
		{"одно число", 1, 100},
		{"ноль чисел", 0, 100},
		{"отрицательное количество", -5, 100},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			want := max(tc.n, 0)

			got, err := collect(randInts(tc.n, tc.maxInt), want)
			if err != nil {
				t.Fatalf("randInts(%d, %d): %v", tc.n, tc.maxInt, err)
			}

			if len(got) != want {
				t.Errorf("получено %d значений, want %d", len(got), want)
			}

			for i, value := range got {
				if value < 0 || value >= tc.maxInt {
					t.Errorf("значение[%d] = %d вне диапазона [0, %d)", i, value, tc.maxInt)
				}
			}
		})
	}
}

func TestRandInts_Unbuffered(t *testing.T) {
	ch := randInts(1, 10)

	if got := cap(ch); got != 0 {
		t.Errorf("cap(канала) = %d, want 0", got)
	}

	if _, err := collect(ch, 1); err != nil {
		t.Fatalf("randInts(1, 10): %v", err)
	}
}

func TestRandInts_ClosesChannel(t *testing.T) {
	ch := randInts(3, 10)

	if _, err := collect(ch, 3); err != nil {
		t.Fatalf("randInts(3, 10): %v", err)
	}

	select {
	case value, ok := <-ch:
		if ok {
			t.Errorf("канал не закрыт после выдачи всех значений: получено %d", value)
		}
	case <-time.After(time.Second):
		t.Fatal("канал не закрыт: чтение блокируется")
	}
}

func TestRandInts_ReturnsDifferentValues(t *testing.T) {
	const (
		n      = 20
		maxInt = 1_000_000
	)

	first, err := collect(randInts(n, maxInt), n)
	if err != nil {
		t.Fatalf("первый генератор: %v", err)
	}

	second, err := collect(randInts(n, maxInt), n)
	if err != nil {
		t.Fatalf("второй генератор: %v", err)
	}

	if slices.Equal(first, second) {
		t.Errorf("два генератора вернули одинаковые последовательности: %v", first)
	}
}
