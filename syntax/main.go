package main

import (
	"fmt"
	"learn_go/syntax/slice"
)

func main() {
	s1 := []int{}
	s2 := []int{1}
	s3 := []int{1, 2, 3, 4, 5}

	r1, _ := slice.DeleteAt(s1, 0)
	fmt.Printf("%v\n", r1)
	r2, _ := slice.DeleteAt(s2, 0)
	fmt.Printf("%v\n", r2)
	r3, _ := slice.DeleteAt(s3, 0)
	r3, _ = slice.DeleteAt(r3, 3)
	fmt.Printf("%v\n", r3)

}
