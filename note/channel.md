```
ch, ok := make(chan int, 2)

// ok可以用于判定channel是否关闭
```



```
sudo php artisan withdrawal_type:add --type=151

sudo chown -R www:www ./storage
```



可以利用第二个返回值来判断channel返回的值是否有效。



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

