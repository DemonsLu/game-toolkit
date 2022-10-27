package main

import (
	"sync"
	"testing"
	"time"
)

func TestWheel(t *testing.T) {
	wg := sync.WaitGroup{}
	wg.Add(1)
	InitTimerManager()
	StartTimerManager()

	AddMultiExecTimer(1000, -1, 5000, func(param any) {
		t.Log("forever")
	})

	AddMultiExecTimer(1000, 10, 2000, func(param any) {
		t.Log("tickA")
	})

	timerB := AddMultiExecTimer(10, 100, 100, func(param any) {
		t.Log("tickB")
	})

	AddTimer(10000, func(param any) {
		AddMultiExecTimer(10, 100, 100, func(param any) {
			t.Log("like but not tickB")
		})
	})

	go func() {
		time.Sleep(3 * time.Second)
		DeleteTimer(timerB)
	}()

	go func() {
		time.Sleep(20 * time.Second)
		wg.Done()
	}()

	wg.Wait()
}
