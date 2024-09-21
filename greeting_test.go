package main

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestHello(t *testing.T) {
	name := "Tony"
	msg, _ := Hello("")

	fmt.Println(1)

	// 与assert包实现相同的断言，但是在测试失败时停止测试执行。

	// 期望没有错误，如果发生错误则停止运行。
	//require.NoError(t, err)
	assert.Equal(t, msg, "hi, "+name)

	fmt.Println(2)
}

func TestHelloEmpty(t *testing.T) {
	msg, err := Hello("")
	if msg != "" || err == nil {
		t.Fatalf(`Hello("") = %q, %v want "", error`, msg, err)
	}
}
