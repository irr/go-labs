package main

import (
	"fmt"
)

type Node struct {
	Data  *string
	Left  *Node
	Right *Node
}

type Tree struct {
	Root *Node
	N    int
}

func newNode(data *string) *Node {
	node := new(Node)
	node.Data = data
	return node
}

func (t *Tree) minValue(node *Node) *Node {
	p := node
	for p.Left != nil {
		p = p.Left
	}
	return p
}

func (t *Tree) Remove(data string) bool {
	return (t.remove(t.Root, &data) != nil)
}

func (t *Tree) remove(node *Node, data *string) *Node {
	if node == nil {
		return nil
	}
	if *data < *node.Data {
		node.Left = t.remove(node.Left, data)
	} else if *data > *node.Data {
		node.Right = t.remove(node.Right, data)
	} else {
		if node.Left == nil {
			if node == t.Root {
				t.Root = node.Right
			}
			t.N--
			return node.Right
		} else if node.Right == nil {
			if node == t.Root {
				t.Root = node.Left
			}
			t.N--
			return node.Left
		}
		p := t.minValue(node.Right)
		node.Data = p.Data
		node.Right = t.remove(node.Right, p.Data)
	}
	return node
}

func (t *Tree) Search(data string) *Node {
	return t.lookup(t.Root, &data)
}

func (t *Tree) lookup(node *Node, data *string) *Node {
	if node == nil {
		return nil
	}
	if *data == *node.Data {
		return node
	} else {
		if *data < *node.Data {
			return t.lookup(node.Left, data)
		} else {
			return t.lookup(node.Right, data)
		}
	}
}

func (t *Tree) Add(data string) *Node {
	node := t.insert(t.Root, &data)
	if t.Root == nil {
		t.Root = node
	}
	t.N++
	return node
}

func (t *Tree) insert(node *Node, data *string) *Node {
	if node == nil {
		return newNode(data)
	} else {
		if *data < *node.Data {
			node.Left = t.insert(node.Left, data)
		} else {
			node.Right = t.insert(node.Right, data)
		}
		return node
	}
}

func (t *Tree) Dump() {
	if t.Root == nil {
		fmt.Println("<empty>")
	}
	fmt.Printf("Tree: %d\n", t.N)
	dumpTree(t.Root)
}

func dumpTree(node *Node) {
	if node == nil {
		return
	}
	dumpTree(node.Left)
	fmt.Printf("\t%#v\n", *node.Data)
	dumpTree(node.Right)
}

func NewTree() *Tree {
	tree := new(Tree)
	tree.Root = nil
	return tree
}

func main() {
	tree := NewTree()
	for _, v := range []string{"ale", "irr", "luma", "lara", "babi"} {
		tree.Add(v)
	}
	tree.Dump()
	node := tree.Search("luma")
	fmt.Printf("Search: luma\n\t%#v[%v]\n", node, *node.Data)
	key := "luma"
	fmt.Printf("Remove: %s\n\t%#v\n", key, tree.Remove(key))
	tree.Dump()
}
