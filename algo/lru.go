package main

import (
	"fmt"
	"strconv"
)

type Node struct {
	Prev *Node
	Next *Node
	Data *string
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

type LRUData struct {
	Node  *Node
	Value *string
}

type LRU struct {
	N    int
	Link *List
	Hash map[string]LRUData
}

func (l *List) Push(d *string) *Node {
	node := new(Node)
	node.Data = d
	if l.Last == nil {
		l.First = node
		l.Last = node
	} else {
		l.Last.Next = node
		node.Prev = l.Last
		l.Last = node
	}
	l.N++
	return node
}

func (l *List) Remove(node *Node) *string {
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

func (l *List) Shift() *string {
	if l.First == nil {
		return nil
	}
	data := l.First.Data
	l.Remove(l.First)
	return data
}

func (l *List) Pop() *string {
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
		fmt.Printf("%02d:%#v[%v]\n", i, p, *p.Data)
		p = p.Next
		i++
	}
}

func NewLRU(n int) *LRU {
	lru := new(LRU)
	lru.N = n
	lru.Link = NewList()
	lru.Hash = make(map[string]LRUData)
	return lru
}

func (lru *LRU) Delete(k string) {
	data, ok := lru.Hash[k]
	if ok {
		lru.Link.Remove(data.Node)
		delete(lru.Hash, *data.Value)
	}
}

func (lru *LRU) Add(k, v string) {
	if lru.Link.N == lru.N {
		id := lru.Link.Shift()
		fmt.Printf("removing %s from linked list (limit reached)\n", *id)
		delete(lru.Hash, *id)
	}
	lru.Delete(k)
	node := lru.Link.Push(&k)
	lru.Hash[k] = LRUData{Node: node, Value: &v}
}

func (lru *LRU) Dump() {
	fmt.Printf("LRU(List):%d\n", lru.Link.N)
	lru.Link.Dump()
	fmt.Println("LRU(Map):")
	if len(lru.Hash) > 0 {
		for k, v := range lru.Hash {
			fmt.Printf("%#v: %#v\n", k, *v.Value)
		}
	} else {
		fmt.Println("<empty>")
	}
}

func testList() {
	var m *Node
	d := NewList()
	max := 10
	for i := 1; i <= max; i++ {
		c := strconv.Itoa(i)
		n := d.Push(&c)
		if i == max/2 {
			m = n
		}
	}
	d.Dump()
	fmt.Printf("Pop: %#v\n", *d.Pop())
	d.Dump()
	fmt.Printf("Shift: %#v\n", *d.Shift())
	d.Dump()
	fmt.Printf("Remove: %#v\n", *d.Remove(m))
	d.Dump()
}

func testLRU() {
	max := 5
	lru := NewLRU(max)
	for i := 1; i <= max; i++ {
		c := strconv.Itoa(i)
		lru.Add(c, "Data"+c)
	}
	lru.Dump()
	fmt.Println("adding 6 to LRU")
	lru.Add("6", "Data6")
	lru.Dump()
	fmt.Println("adding 2 to LRU")
	lru.Add("2", "Data2")
	lru.Dump()
	fmt.Println("adding 7 to LRU")
	lru.Add("7", "Data7")
	lru.Dump()
}

func main() {
	testLRU()
}
