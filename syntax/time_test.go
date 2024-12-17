package main

import (
	"context"
	"testing"
	"time"
)

func TestTicker(t *testing.T) {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	for {
		select {
		case <-ctx.Done():
			t.Log("超时了，或者被取消了")
			goto end
		case now := <-ticker.C:
			t.Log(now.Unix())
		}
	}
end:
}

func TestDefer(t *testing.T) {
	defer func() {
		t.Log("task1")
	}()
	defer func() {
		t.Log("task2")
	}()
	defer func() {
		t.Log("task3")
	}()
}
