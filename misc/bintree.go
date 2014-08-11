package main

type Item interface {
	Less(than Item) bool
	Equal(to Item) bool
}

type BinaryTree struct {
	node  Item
	left  *BinaryTree
	right *BinaryTree
}

func New() *BinaryTree {
	tree := &BinaryTree{}
	tree.node = nil
	return tree
}

func (tree *BinaryTree) Search(value Item) *BinaryTree {
	if tree.node == nil {
		return nil
	}

	if tree.node.Equal(value) {
		return tree
	} else {
		if value.Less(tree.node) == true {
			t := tree.left.Search(value)
			return t
		} else {
			t := tree.right.Search(value)
			return t
		}
	}
}

func (tree *BinaryTree) Insert(value Item) {
	if tree.node == nil {
		tree.node = value
		tree.right = New()
		tree.left = New()
		return
	} else {
		if value.Less(tree.node) == true {
			tree.left.Insert(value)
		} else {
			tree.right.Insert(value)
		}
	}
}

type Int int

func (this Int) Less(than Item) bool {
	return this < than.(Int)
}

func (this Int) Equal(to Item) bool {
	return this == to.(Int)
}

func main() {
	tree := New()

	tree.Insert(Int(1))
	tree.Insert(Int(2))
	tree.Insert(Int(3))
}
