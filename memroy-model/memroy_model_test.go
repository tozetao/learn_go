package memroy_model

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

/*
参考：
https://go.dev/ref/mem
https://jasonkayzk.github.io/2022/10/26/%E3%80%90%E7%BF%BB%E8%AF%91%E3%80%91Go%E5%86%85%E5%AD%98%E6%A8%A1%E5%9E%8B/

内存模型：讨论的是内存一致性模型。

happends before
	如果a1操作发生于a2操作之前，就可以说a1 happends before a2，此时a1操作的变量结果对a2是可见的。
	注：指令重排不会影响happends before

a := 1
b := 2
书写顺序是a先赋值，然后b再赋值，但是在转换成机器代码时，可能是b先赋值，这就是指令重排。

go的哪些行为是保证happends before的？

与goroutine相关的
 1. goroutine的创建happends before其执行
    个人理解是goroutine的创建顺序可以按照代码的书写顺序来创建的。

 2. goroutine的完成不保证happends bofore任何代码
    虽然goroutine的创建顺序可以按照代码的书写顺序来创建，但是不同goroutine中的代码谁先执行谁后执行是不保证的。

*/

func TestHappendsBefore1(t *testing.T) {
	var a, b int

	// a和b俩个分支都可能输出，go不保证goroutine代码执行的happends before
	go func() {
		if b == 2 {
			print(fmt.Sprintf("a: %d\n", a))
		} else {
			print(fmt.Sprintf("b: %d\n", b))
		}
	}()

	go func() {
		a = 1
		b = 2
	}()

	time.Sleep(time.Second * 2)
}

// 与channel相关的
// 1. 对于无缓存通道，读会先于写发生。
// 2. 对于有缓存通道，写会咸鱼读发生。
func TestChannel(t *testing.T) {

}

/*
锁相关
对于任意的sync.Mutex或者sync.RWMutex，n次Unlock调用happends before m次Lock()调用，其中n < m
*/
func TestLocker(t *testing.T) {
	var l sync.Mutex
	var a string

	l.Lock()
	go func() {
		a = "hello world"
		l.Unlock()
	}()
	l.Lock()
	println(a)
}
