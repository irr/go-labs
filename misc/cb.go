package main

import (
	"errors"
	"fmt"
)

type CAPArray struct {
	N int
	a []interface{}
}

func (c *CAPArray) init() {
	if c.a == nil {
		c.a = make([]interface{}, 0, c.N+1)
	}
}

func (c *CAPArray) Fill(e interface{}) {
	for i := 0; i < c.N; i++ {
		c.Push(e)
	}
}

func (c *CAPArray) Geth(i int) (interface{}, error) {
	if len(c.a) != c.N {
		return nil, errors.New("CAPArray must be full")
	}
	if i <= 0 || i > c.N {
		return nil, errors.New("index out of range")
	}
	return c.a[c.N-i], nil
}

func (c *CAPArray) Seth(i int, e interface{}) error {
	if len(c.a) != c.N {
		return errors.New("CAPArray must be full")
	}
	c.a[c.N-i] = e
	return nil
}

func (c *CAPArray) Push(e interface{}) {
	c.init()
	c.a = append(c.a, e)
	if len(c.a) > c.N {
		c.a = c.a[1:len(c.a)]
	}
}

func (c *CAPArray) Pop() (e interface{}) {
	if c.a == nil {
		return
	}
	e = c.a[len(c.a)-1]
	c.a = c.a[:len(c.a)-1]
	return e
}

func (c *CAPArray) Unshift(e interface{}) {
	c.init()
	c.a = append(c.a, e)
	copy(c.a[1:], c.a[0:])
	c.a[0] = e
	if len(c.a) > c.N {
		c.a = c.a[0 : len(c.a)-1]
	}
}

func (c *CAPArray) Shift() (e interface{}) {
	if c.a == nil {
		return
	}
	e = c.a[0]
	c.a = c.a[1:]
	return e
}

func (c *CAPArray) Dump(f string) {
	if len(c.a) > 0 {
		fmt.Printf(fmt.Sprintf("n=%s first=%s last=%s { ", f, f, f), len(c.a), c.a[0], c.a[len(c.a)-1])
		for _, v := range c.a {
			fmt.Printf(f, v)
		}
		fmt.Println("}")
	} else {
		fmt.Println("{}")
	}
}

func main() {
	s2 := CAPArray{N: 20}
	fmt.Printf("\nCAPArray(%#v)\n====================================================\n", s2)
	fmt.Println("Push...")
	for i := 0; i < 100; i++ {
		fmt.Printf("e=%02d ", i)
		s2.Push(i)
		s2.Dump("%02d ")
	}
	fmt.Println("Pop...")
	for i := 0; i < 5; i++ {
		fmt.Printf("e=%02d ", s2.Pop())
		s2.Dump("%02d ")
	}
	fmt.Println("Shift...")
	for i := 0; i < 5; i++ {
		fmt.Printf("e=%02d ", s2.Shift())
		s2.Dump("%02d ")
	}
	fmt.Println("UnShift...")
	for i := 84; i > 65; i-- {
		fmt.Printf("e=%02d ", i)
		s2.Unshift(i)
		s2.Dump("%02d ")
	}

	s3 := CAPArray{N: 10}
	fmt.Printf("\nCAPArray(%#v)\n====================================================\n", s3)
	fmt.Printf("Pop...   (%#v)\ne=%#v\n", s3, s3.Pop())
	fmt.Printf("Shift... (%#v)\ne=%#v\n", s3, s3.Shift())

	fmt.Printf("\nTesting Geth and Seth...\n")
	fmt.Printf("\nCAPArray(%#v)\n====================================================\n", s3)
	_, err := s3.Geth(10)
	if err != nil {
		fmt.Println(err)
	}

	for i := 0; i < 10; i++ {
		s3.Push(i)
	}

	s3.Dump("%02d ")
	fmt.Println()

	get := func(i int) {
		e, err := s3.Geth(i)
		if err != nil {
			fmt.Println("get(", i, ") ", err)
		} else {
			fmt.Printf("s3[%d]=%v\n", i, e)
		}
	}

	set := func(i int, e interface{}) {
		err := s3.Seth(i, e)
		if err != nil {
			fmt.Println("set(", i, ") ", err)
		}
	}

	get(1)
	get(2)
	get(10)
	get(0)
	get(11)

	set(1, 11)
	get(1)
	set(3, 33)
	get(3)
	set(9, 99)

	fmt.Println()

	s3.Dump("%02d ")
}
