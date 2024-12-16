package main

import (
	"fmt"
	"github.com/robfig/cron/v3"
	"testing"
	"time"
)

// job不属于ddd的领域范畴。
// service的调用方可以是web、rpc、command、job
// 我们先实现service，再实现job

func TestCron(t *testing.T) {
	//ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	//defer cancel()

	c := cron.New()

	showStrJob := &ShowStr{}
	c.AddJob("@every 3s", showStrJob)

	c.Start()
	time.Sleep(time.Second * 10)

	stop := c.Stop()
	t.Log("准备结束了")
	<-stop.Done()
	t.Log("结束了")

	//select {
	//case <-ctx.Done():
	//	t.Log("超时，程序结束了")
	//}
	//t.Log("结束了")
}

type ShowStr struct {
	ID int
}

func (s *ShowStr) Run() {
	fmt.Println("运行showstr任务")
	time.Sleep(time.Second * 10)
	fmt.Println("任务结束")
}

func TestMap(t *testing.T) {
	m := make(map[int]ShowStr, 10)

	m[1] = ShowStr{ID: 10}
	m[2] = ShowStr{ID: 20}

	job1 := m[1]
	fmt.Printf("joob1: %v\n", job1.ID)

	job3 := m[3]
	fmt.Printf("job3: %v\n", job3.ID)

	job4, exists := m[4]
	fmt.Printf("job3: %v, exists: %v\n", job4.ID, exists)
}
