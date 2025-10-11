package main

import "fmt"

func modify(m map[string]int) {
	m["x"] = 100
}

func reassign(m map[string]int) {
	m["y"] = 200
}

func main() {
	data := map[string]int{"a": 1}

	modify(data)
	fmt.Println("After modify:", data) // shows x:100 added

	reassign(data)
	fmt.Println("After reassign:", data) // unchanged
}
