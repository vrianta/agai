package pac

import "fmt"

type New struct {
	Test string
}

func (n New) Print() {
	fmt.Printf("This is a pacman print function, Test : %s\n", n.Test)
}

func (n *New) Init(test string) {
	n.Test = test
}
