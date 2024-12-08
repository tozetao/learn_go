package channel

import (
	"fmt"
	"testing"
	"time"
)

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
