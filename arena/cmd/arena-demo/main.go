package main

import (
	"arena"
	"fmt"
)

// Packet models a network payload allocated inside an arena.
type Packet struct {
	ID   int
	Data []byte
}

func main() {
	a := arena.NewArena()
	defer a.Free()

	p := arena.New[Packet](a)
	p.ID = 1
	p.Data = []byte("payload")

	fmt.Printf("Packet: %+v\n", *p)
}
