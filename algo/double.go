package main

import (
	"fmt"
	"strconv"
)

type Content struct {
	Id string
}

type Node struct {
	Prev *Node
	Next *Node
	Data *Content
}

type List struct {
	N     int
	First *Node
	Last  *Node
}

func NewList() *List {
	list := new(List)
	return list
}

func (l *List) Push(d *Content) *Node {
	node := new(Node)
	node.Data = d
	if l.Last == nil {
		l.First = node
		l.Last = node
		l.N++
	} else {
		l.Last.Next = node
		node.Prev = l.Last
		l.Last = node
	}
	return node
}

func (l *List) Remove(node *Node) *Content {
	if node == nil {
		return nil
	}
	if node == l.First && node == l.Last {
		l.First = nil
		l.Last = nil
	} else if node == l.First {
		l.First = node.Next
		l.First.Prev = nil
	} else if node == l.Last {
		l.Last = node.Prev
		l.Last.Next = nil
	} else {
		after := node.Next
		before := node.Prev
		after.Prev = before
		before.Next = after
	}
	l.N--
	return node.Data
}

func (l *List) Shift() *Content {
	if l.First == nil {
		return nil
	}
	data := l.First.Data
	l.Remove(l.First)
	return data
}

func (l *List) Pop() *Content {
	if l.Last == nil {
		return nil
	}
	data := l.Last.Data
	l.Remove(l.Last)
	return data
}

func (l *List) Dump() {
	if l.First == nil {
		fmt.Println("<empty>")
	}
	p := l.First
	i := 1
	for p != nil {
		fmt.Printf("%02d:%#v[%s]\n", i, p, p.Data.Id)
		p = p.Next
		i++
	}
}

func main() {
	var m *Node
	d := NewList()
	max := 10
	for i := 1; i <= max; i++ {
		c := new(Content)
		c.Id = strconv.Itoa(i)
		n := d.Push(c)
		if i == max/2 {
			m = n
		}
	}
	d.Dump()
	fmt.Printf("Pop: %#v\n", d.Pop())
	d.Dump()
	fmt.Printf("Shift: %#v\n", d.Shift())
	d.Dump()
	fmt.Printf("Remove: %#v\n", d.Remove(m))
	d.Dump()
}
