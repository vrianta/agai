package pac

import (
	"fmt"
)

type New struct {
	Test string
	N    *New
}

type PacInterface interface {
	Print()
	Init(test string)
	Clone() PacInterface
}

func (n *New) Print() {
	fmt.Printf("This is a pacman print function, Test : %s\n", n.Test)
}

func (n *New) Init(test string) {
	n.Test = test
}

func (n *New) Clone() PacInterface {
	// 1. Create a shallow copy of the struct value
	t := *n

	// 2. Clear any nested pointers if necessary for deep copy
	// (If you want a true deep clone and not just a shallow copy)
	// t.N = nil

	// 3. FIX: Return the ADDRESS of the copied value.
	// The type *New fully implements PacInterface.
	return &t
}
