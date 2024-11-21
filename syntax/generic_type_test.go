package main

import (
	"fmt"
	"testing"
)

// https://segmentfault.com/a/1190000041634906

type Woo[T int | string] int

func TestGenericType(t *testing.T) {
	var a Woo[int] = 123
	var b Woo[string] = 100
	fmt.Println(a, b)
}
