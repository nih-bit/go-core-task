package main

import (
	"context"
	"fmt"
	"slices"
	"testing"
	"time"
)

// collect вычитывает канал до закрытия.
// Таймаут не даёт тесту зависнуть, если стадия не закроет выходной канал.
func collect(ch <-chan float64, limit int) ([]float64, error) {
	timeout := time.After(2 * time.Second)
	got := make([]float64, 0, limit)

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

func TestCubeFloat64(t *testing.T) {
	tests := []struct {
		name  string
		input uint8
		want  float64
	}{
		{"ноль", 0, 0},
		{"единица", 1, 1},
		{"тройка", 3, 27},
		{"пять", 5, 125},
		{"сто", 100, 1_000_000},
		{"максимум uint8", 255, 16_581_375},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			got, err := collect(cubeFloat64(ctx, producer(ctx, tc.input)), 1)
			if err != nil {
				t.Fatalf("cubeFloat64(%d): %v", tc.input, err)
			}

			if len(got) != 1 {
				t.Fatalf("получено %d значений, want 1: %v", len(got), got)
			}

			if got[0] != tc.want {
				t.Errorf("куб %d = %v, want %v", tc.input, got[0], tc.want)
			}
		})
	}
}

func TestCubeFloat64_PreservesOrder(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	input := []uint8{7, 3, 0, 255, 1}
	want := []float64{343, 27, 0, 16_581_375, 1}

	got, err := collect(cubeFloat64(ctx, producer(ctx, input...)), len(input))
	if err != nil {
		t.Fatalf("cubeFloat64: %v", err)
	}

	if !slices.Equal(got, want) {
		t.Errorf("cubeFloat64(%v) = %v, want %v", input, got, want)
	}
}

func TestCubeFloat64_EmptyInput(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	got, err := collect(cubeFloat64(ctx, producer(ctx)), 0)
	if err != nil {
		t.Fatalf("cubeFloat64: %v", err)
	}

	if len(got) != 0 {
		t.Errorf("cubeFloat64 на пустом входе вернул %v, want пусто", got)
	}
}

func TestCubeFloat64_ClosesOut(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	out := cubeFloat64(ctx, producer(ctx, 2))

	if _, err := collect(out, 1); err != nil {
		t.Fatalf("cubeFloat64: %v", err)
	}

	select {
	case value, ok := <-out:
		if ok {
			t.Errorf("канал не закрыт после закрытия входа: получено %v", value)
		}
	case <-time.After(time.Second):
		t.Fatal("канал не закрыт: чтение блокируется")
	}
}

func TestCubeFloat64_Cancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())

	never := make(chan uint8)

	out := cubeFloat64(ctx, never)

	cancel()

	select {
	case value, ok := <-out:
		if ok {
			t.Errorf("получено значение %v, ожидалось закрытие", value)
		}
	case <-time.After(time.Second):
		t.Fatal("out не закрыт после отмены: стадия висит на приёме из in")
	}
}

func TestProducer_Cancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())

	ch := producer(ctx, 1, 2, 3)

	cancel()

	select {
	case value, ok := <-ch:
		if ok {
			t.Errorf("получено значение %v, ожидалось закрытие", value)
		}
	case <-time.After(time.Second):
		t.Fatal("producer не завершился после отмены: висит на отправке")
	}
}
