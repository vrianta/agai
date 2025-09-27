package main

import "github.com/vrianta/agai/tests/pac"

type pacInterface interface {
	Print()
	Init(test string)
}

func main() {
	p1 := pac.New{}
	p2 := p1

	p1.Init("p1")
	p2.Init("p2")

	p1.Print()
	p2.Print()

}

func createNewPac[T pacInterface]() func() pacInterface {
	return func() pacInterface {
		var t T
		return t
	}
}
