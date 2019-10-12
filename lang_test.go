package bbrpc

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

// golang 测试

func TestFailInGoroutine(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()

		t.Fatal("fatal in goroutine")
	}()
	wg.Wait()
	t.Log("done")
}

func TestRunningRoutine(t *testing.T) {
	go func() {
		tk := time.NewTicker(time.Second)
		for v := range tk.C {
			fmt.Println(v)
		}
	}()
	fmt.Println("end")
}
