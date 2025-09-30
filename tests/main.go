package main

import "github.com/vrianta/agai/tests/pac"

type t2Pac struct {
	pac.New
}

func (p2 *t2Pac) Print() {
	println("t2 Print")
}

func main() {
	p1 := t2Pac{}
	fp := getobj_create_func[t2Pac]()

	p2 := fp()

	p2.Init("testing")
	p2.Print()
	test_interface_cloning(&p1)

}

// THE FINAL, CANONICAL SOLUTION
func getobj_create_func[T any, PT interface {
	*T
	pac.PacInterface
}]() func() pac.PacInterface {
	return func() pac.PacInterface {
		// 1. `new(T)` creates a value of type *T.
		// 2. We declare a variable `p` of type `PT`.
		// 3. The compiler knows `PT`'s underlying type is `*T`, so the assignment is valid.
		// 4. Critically, the compiler also knows `PT` satisfies the interface.
		var p PT = new(T)

		// 5. By returning `p`, we are returning a value whose type (`PT`) is
		//    guaranteed by the constraint to be assignable to the interface.
		return p
	}
}

func test_interface_cloning(p pac.PacInterface) {

	p2 := p.Clone()
	// p3 := pac.New{}

	p.Init("p1")
	p2.Init("p2")
	// p3.Init("p3")
	// p2.N = &p3

	p.Print()
	p2.Print()
	// p1.N.Print()

}
