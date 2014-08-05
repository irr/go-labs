package main

import (
	"code.google.com/p/go-tour/tree"
	"fmt"
)

func Walk(t *tree.Tree, ch chan int) {
	if t.Left != nil {
		Walk(t.Left, ch)
	}
	ch <- t.Value
	if t.Right != nil {
		Walk(t.Right, ch)
	}
}

func Same(t1, t2 *tree.Tree) bool {
	ch1 := make(chan int)
	ch2 := make(chan int)

	m := LeafList{make(map[int]bool)}

	go Walk(t1, ch1)
	go Walk(t2, ch2)

	for i := 0; i < 10; i++ {
		m.DeleteOrAdd(<-ch1)
		m.DeleteOrAdd(<-ch2)
	}

	return len(m.List) == 0

}

type LeafList struct {
	List map[int]bool
}

func (m *LeafList) DeleteOrAdd(v int) {
	_, ok := m.List[v]
	if !ok {
		m.List[v] = true
	} else {
		delete(m.List, v)
	}
}

func main() {
	fmt.Println(Same(tree.New(1), tree.New(1)))
	fmt.Println(Same(tree.New(1), tree.New(2)))
}
