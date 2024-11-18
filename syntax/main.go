package main

import (
	"context"
	"fmt"
	"time"
)

func main() {
	//s1 := []int{}
	//s2 := []int{1}
	//s3 := []int{1, 2, 3, 4, 5}
	//
	//r1, _ := slice.DeleteAt(s1, 0)
	//fmt.Printf("%v\n", r1)
	//r2, _ := slice.DeleteAt(s2, 0)
	//fmt.Printf("%v\n", r2)
	//r3, _ := slice.DeleteAt(s3, 0)
	//r3, _ = slice.DeleteAt(r3, 3)
	//fmt.Printf("%v\n", r3)

	TestContext()

	time.Sleep(20 * time.Second)
}

func TestContext() {
	ctx, cancel := context.WithCancel(context.Background())

	go func(ctx context.Context) {
		select {
		case <-ctx.Done():
			fmt.Println("任务被取消了")
		case <-time.After(time.Second * 15):
			fmt.Println("任务完成")
		}
	}(ctx)

	fmt.Println("睡眠5秒")
	time.Sleep(time.Second * 5)
	fmt.Println("取消任务")
	cancel()
}
