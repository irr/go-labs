package main

import (
	"fmt"
)

type BinHeap struct {
	list []int
	size int
}

func (b *BinHeap) percUp(i int) {
	for ; (i / 2) > 0; i = i / 2 {
		if b.list[i] < b.list[i/2] {
			b.list[i], b.list[i/2] = b.list[i/2], b.list[i]
		}
	}
}

func (b *BinHeap) percDown(i int) {
	for i*2 <= b.size {
		mc := b.minChild(i)
		if b.list[i] > b.list[mc] {
			b.list[i], b.list[mc] = b.list[mc], b.list[i]
		}
		i = mc
	}
}

func (b *BinHeap) minChild(i int) int {
	if (i*2 + 1) > b.size {
		return i * 2
	} else {
		if b.list[i*2] < b.list[i*2+1] {
			return i * 2
		} else {
			return i*2 + 1
		}
	}
}

func NewBinHeap() *BinHeap {
	b := new(BinHeap)
	b.size = 0
	b.list = []int{0}
	return b
}

func (b *BinHeap) insert(k int) {
	b.list = append(b.list, k)
	b.size++
	b.percUp(b.size)
}

func (b *BinHeap) delmin() (int, bool) {
	if len(b.list) < 2 {
		return 0, false
	}
	r := b.list[1]
	b.list[1] = b.list[b.size]
	b.size--
	b.list = b.list[:len(b.list)-1]
	b.percDown(1)
	return r, true
}

func BuildHeap(alist []int) *BinHeap {
	if alist == nil || len(alist) == 0 {
		return NewBinHeap()
	}
	b := NewBinHeap()
	i := len(alist) / 2
	b.size = len(alist)
	b.list = append(alist, 0)
	copy(b.list[1:], alist[0:])
	b.list[0] = 0
	for i > 0 {
		b.percDown(i)
		i--
	}
	return b
}

func main() {
	b := NewBinHeap()
	fmt.Printf("%#v\n", b)

	v, ok := b.delmin()
	fmt.Printf("%#v[%v]\n", v, ok)

	for _, v := range []int{33, 17, 27, 18, 14, 9, 19, 21, 11, 5, 7} {
		b.insert(v)
		fmt.Printf("%#v\n", b)
	}

	for i := 0; i < 3; i++ {
		v, ok := b.delmin()
		fmt.Printf("%#v[%v]\n", v, ok)
		fmt.Printf("%#v\n", b)
	}

	for _, v := range []int{1, 100} {
		b.insert(v)
		fmt.Printf("%#v\n", b)
	}

	bh := BuildHeap([]int{33, 17, 27, 18, 14, 9, 19, 21, 11, 5, 7})
	fmt.Printf("%#v\n", bh)
	for i := 0; i < 11; i++ {
		v, ok := bh.delmin()
		fmt.Printf("%#v[%v]\n", v, ok)
		fmt.Printf("%#v\n", bh)
	}
}
