package defer_demo

func DeferV1() {
	i := 0
	i = 1
	defer func() {
		println(i)
	}()
	i = 2
	defer func() {
		println(i)
	}()
}

func DeferV2() {
	i := 0
	defer func(i int) {
		println(i)
	}(i)
	i = 1
}

func DeferClosureLoopV1() {
	for i := 0; i < 10; i++ {
		defer func() {
			println(i)
		}()
	}
}

func DeferClosureLoopV2() {
	for i := 0; i < 10; i++ {
		defer func(val int) {
			println(val)
		}(i)
	}
}

func DeferClosureLoopV3() {
	for i := 0; i < 10; i++ {
		j := i
		defer func() {
			println(j)
		}()
	}
}
