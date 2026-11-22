package main

import (
	"fmt"
	"slices"
	"testing"
	"time"
)

// создаёт канал с заданными значениями и закрывает его после отправки.
func gen(values ...int) <-chan int {
	ch := make(chan int)

	go func() {
		defer close(ch)

		for _, value := range values {
			ch <- value
		}
	}()

	return ch
}

// вычитывает канал до закрытия, таймаут не даёт тесту зависнуть
func drain(ch <-chan int, limit int) ([]int, error) {
	timeout := time.After(2 * time.Second)
	got := make([]int, 0, limit)

	for {
		select {
		case value, ok := <-ch:
			if !ok {
				return got, nil
			}

			got = append(got, value)

			if len(got) > limit {
				return got, fmt.Errorf("получено больше %d значений — канал не закрывается", limit)
			}
		case <-timeout:
			return got, fmt.Errorf("таймаут: получено %d значений из %d", len(got), limit)
		}
	}
}

func TestMerge(t *testing.T) {
	tests := []struct {
		name   string
		inputs [][]int
		want   []int
	}{
		{"три канала", [][]int{{1, 2}, {3, 4, 5}, {6, 7, 8}}, []int{1, 2, 3, 4, 5, 6, 7, 8}},
		{"один канал", [][]int{{1, 2, 3}}, []int{1, 2, 3}},
		{"ни одного канала", nil, []int{}},
		{"пустые каналы", [][]int{{}, {}}, []int{}},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			chs := make([]<-chan int, len(tc.inputs))
			for i, values := range tc.inputs {
				chs[i] = gen(values...)
			}

			got, err := drain(merge(chs...), len(tc.want))
			if err != nil {
				t.Fatalf("merge: %v", err)
			}

			slices.Sort(got)

			if !slices.Equal(got, tc.want) {
				t.Errorf("merge вернул %v, want %v", got, tc.want)
			}
		})
	}
}

func TestMerge_ClosesOut(t *testing.T) {
	out := merge(gen(1, 2, 3))

	if _, err := drain(out, 3); err != nil {
		t.Fatalf("merge: %v", err)
	}

	select {
	case value, ok := <-out:
		if ok {
			t.Errorf("канал не закрыт после закрытия входов: получено %d", value)
		}
	case <-time.After(time.Second):
		t.Fatal("канал не закрыт: чтение блокируется")
	}
}

func TestMerge_NilChannel(t *testing.T) {
	t.Run("только nil", func(t *testing.T) {
		t.Parallel()

		got, err := drain(merge(nil), 0)
		if err != nil {
			t.Fatalf("merge(nil): %v", err)
		}

		if len(got) != 0 {
			t.Errorf("merge(nil) вернул %v, want пусто", got)
		}
	})

	t.Run("nil среди каналов", func(t *testing.T) {
		t.Parallel()

		got, err := drain(merge(gen(1, 2), nil, gen(3)), 3)
		if err != nil {
			t.Fatalf("merge с nil: %v", err)
		}

		slices.Sort(got)

		if want := []int{1, 2, 3}; !slices.Equal(got, want) {
			t.Errorf("merge(1,2 | nil | 3) = %v, want %v", got, want)
		}
	})
}
