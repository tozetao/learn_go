package main

import (
	"context"
	"fmt"
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

func TestWithValue(t *testing.T) {
	ctx := context.WithValue(context.Background(), "key1", "value1")
	value := ctx.Value("key1")
	t.Log(value)

	ctx2 := context.WithValue(ctx, "key2", "value2")
	t.Log(ctx2.Value("key1"))
	t.Log(ctx2.Value("key2"))

	ctx = context.WithValue(ctx, "key1", "hi,demo")
	t.Log(ctx.Value("key1"))
}

func TestInterface(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer func() {
		cancel()
		t.Log(ctx.Err())
	}()

	select {
	case <-ctx.Done():
		t.Log(fmt.Sprintf("ctx.Done, err:%v", ctx.Err()))
	}
}
