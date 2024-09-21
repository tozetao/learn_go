package control

import "fmt"

type User struct {
	Name string
}

func LoopBug() {
	users := []User{
		{
			Name: "小红",
		},
		{
			Name: "王芳",
		},
	}

	m := make(map[string]*User, 2)

	for _, u := range users {
		m[u.Name] = &u
	}

	for k, v := range m {
		fmt.Printf("name: %s, user: %v\n", k, v)
	}
}
