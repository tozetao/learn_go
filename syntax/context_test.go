package main

import (
	"context"
	"testing"
	"time"
)

func TestCtx(t *testing.T) {
	pc, cancel := context.WithCancel(context.Background())

	ctx, childCancel := context.WithCancel(pc)

	go func() {
		select {
		case <-ctx.Done():
			t.Logf("ctx.Done, err:%v", ctx.Err())
		}
	}()

	go func() {
		t.Log("等待关闭父级context")
		time.Sleep(5 * time.Second)
		t.Log("关闭父级context")
		cancel()
	}()

	go func() {
		t.Log("等待关闭子context")
		time.Sleep(10 * time.Second)
		t.Log("关闭子context")
		childCancel()
	}()

	time.Sleep(time.Second * 20)
}
