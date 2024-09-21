package component

type Inner struct {
	Name string
}

func (inner *Inner) DoSomething() {
	println("do something in inner.")
}

type Outer struct {
	Inner
}

func (outer *Outer) DoSomething() {
	println("do someghint in outer.")
}

type Outerv1 struct {
	Inner
	Outer
}
