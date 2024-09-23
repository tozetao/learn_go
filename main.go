package main

import (
	"fmt"
	"github.com/google/uuid"
	"time"
)

func main() {
	ssid := uuid.New()

	fmt.Printf("%v\n", ssid.String())

	//author := &Author{}
	//fmt.Printf("%#v\n", author.createdAt.UnixMilli())
}

type Author struct {
	name      string
	age       int
	createdAt time.Time
}

func (a *Author) fullName() {
	fmt.Printf("%s is %d years old\n", a.name, a.age)
}

type Post struct {
	title   string
	content string
	Author
}

func (p Post) display() {
	fmt.Printf("title: %s\n", p.title)
	p.fullName()
}
