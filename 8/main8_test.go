package main

import (
	"sync/atomic"
	"testing"
	"time"
)

// assertReturns выполняет f и падает, если она не завершилась за timeout.
func assertReturns(t *testing.T, timeout time.Duration, f func()) {
	t.Helper()

	done := make(chan struct{})

	go func() {
		defer close(done)
		f()
	}()

	select {
	case <-done:
	case <-time.After(timeout):
		t.Fatalf("операция не завершилась за %v", timeout)
	}
}

func TestWaitGroup_WaitsForDone(t *testing.T) {
	var wg WaitGroup
	wg.Add(1)

	waitReturned := make(chan struct{})

	go func() {
		wg.Wait()
		close(waitReturned)
	}()

	select {
	case <-waitReturned:
		t.Fatal("Wait вернулся, хотя счётчик равен 1")
	case <-time.After(100 * time.Millisecond):
	}

	wg.Done()

	select {
	case <-waitReturned:
	case <-time.After(time.Second):
		t.Fatal("Wait не вернулся после Done")
	}
}

func TestWaitGroup_ZeroValue(t *testing.T) {
	var wg WaitGroup

	assertReturns(t, time.Second, wg.Wait)
}

func TestWaitGroup_AddZero(t *testing.T) {
	var wg WaitGroup
	wg.Add(0)

	assertReturns(t, time.Second, wg.Wait)
}

func TestWaitGroup_AddAccumulates(t *testing.T) {
	var wg WaitGroup
	wg.Add(3)
	wg.Add(2)

	waitReturned := make(chan struct{})

	go func() {
		wg.Wait()
		close(waitReturned)
	}()

	for i := 0; i < 4; i++ {
		wg.Done()
	}

	select {
	case <-waitReturned:
		t.Fatal("Wait вернулся, хотя счётчик равен 1")
	case <-time.After(100 * time.Millisecond):
	}

	wg.Done()

	select {
	case <-waitReturned:
	case <-time.After(time.Second):
		t.Fatal("Wait не вернулся после последнего Done")
	}
}

func TestWaitGroup_WaitsForAllGoroutines(t *testing.T) {
	const workers = 100

	var wg WaitGroup
	var finished int64

	wg.Add(workers)

	for i := 0; i < workers; i++ {
		go func() {
			defer wg.Done()

			time.Sleep(time.Millisecond)
			atomic.AddInt64(&finished, 1)
		}()
	}

	wg.Wait()

	if got := atomic.LoadInt64(&finished); got != workers {
		t.Errorf("к моменту возврата Wait завершилось %d горутин из %d", got, workers)
	}
}

func TestWaitGroup_Reuse(t *testing.T) {
	var wg WaitGroup

	for cycle := 1; cycle <= 3; cycle++ {
		wg.Add(3)

		for i := 0; i < 3; i++ {
			go func() {
				defer wg.Done()
			}()
		}

		assertReturns(t, time.Second, wg.Wait)
	}
}

func TestWaitGroup_MultipleWaiters(t *testing.T) {
	const waiters = 5

	var wg WaitGroup
	wg.Add(1)

	returned := make(chan struct{}, waiters)

	for i := 0; i < waiters; i++ {
		go func() {
			wg.Wait()
			returned <- struct{}{}
		}()
	}

	select {
	case <-returned:
		t.Fatal("один из ожидающих вернулся раньше Done")
	case <-time.After(100 * time.Millisecond):
	}

	wg.Done()

	for i := 0; i < waiters; i++ {
		select {
		case <-returned:
		case <-time.After(time.Second):
			t.Fatalf("вернулось %d ожидающих из %d", i, waiters)
		}
	}
}
