channel是Go用于并发的一个内置类型。我的理解channel就是goroutine的通道，goroutine之间通过这个channel进行通讯以避免数据竞争问题。

channel的基本操作包括：

- 创建

  ```go
  ch1 := make(chan int)
  
  // 创建带有容量的channel
  ch2 := make(chan struct{}, 2)
  ```

- 发送数据到channel里面

  ```go
  ch <- data
  ```

- 从 channel 里面读取数据

  ```go
  val := <-ch
  ```

- 关闭channel

  ```go
  close(ch)
  ```

  



**channel的close问题**

当一个channel被关闭之后：

- 向其写入数据会panic
- 再次close会panic

但是被close的channel可以读取数据，会读到nil值。

```go
ch := make(chan int, 1)
ch <- 1

// 可以利用第二个返回值来判断channel返回的值是否有效。
n, ok := <-ch
println(n, ok)	// 1, true

// 关闭channel
close(ch)

n1, ok := <-ch	
println(n, ok)	// 0, false
```

在实践中，谁创建的channel谁来关闭，这可以避免很多channel未能正确关闭的问题。但是如果channel是作为参数传递给别人使用，就要注意数据竞态问题。

```go
type MyStruct struct {
    ch chan int
    once Sync.Once
}

func (m MyStruct) SafeClose() {
    once.Do(function() {
        close(ch)
    })
}
```









**channel的阻塞问题**

在使用channel的时候：

如果接收者读不到消息，就会阻塞。分俩种情况：

- 没有缓存，对面没有发送者的时候就会阻塞。
- 如果有缓存，但是channel中没数据，对面没有发送者的时候就会阻塞。

归根结底就是channel有没有数据，对面有没有发送者，这俩个条件是要同时成立就会阻塞。

```go
// 接收者阻塞
func TestChan(t *testing.T) {
	ch := make(chan int)
	go func() {
		for {
			fmt.Println("send a message to ch.")
			ch <- 1
			fmt.Println("start sleep")
			time.Sleep(time.Second * 5)
		}
	}()
	time.Sleep(time.Minute)
}
```



如果发送者写不了消息就会阻塞，也分为俩种情况：

- 如果channel没缓存，且对面没接收者等待接收数据，会阻塞
- 如果channel有缓存，但是缓存满了，且对面没有接收者等待接收数据，会阻塞。







**range循环**

在实践中，常常需要持续的从channel中读取数据，这种可以使用range循环，当channel被关闭了range循环会自动退出。

```go
ch := make(chan int)

go func(){
    for i := 0; i < 10; i++ {
        ch <- i
    }
}()

for val := range ch {
    println(val)
}
println("发送完毕")
```



**channel与select**

在go中还有一个类似switch-case的结构：select-case结构，用于控制从不同的channel中读写数据。

```go
select {
    case ch1 <- val:	// 该分支是写入
    case val := ch2:	// 该分支是读取
    default:
}
```

- 每一个case都可以是读取channel或者写入channel
- 没有default分支时，select会阻塞，直到任何一个case执行成功。
- 如果多个分支满足要求，随机选择一个case执行。
- 当有default分支时，且所有case都阻塞，select将会执行default分支。









