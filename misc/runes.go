package main

import "fmt"

func main() {
	s := "語本日"
	fmt.Printf("[%s] {%#v} len: %d, %d\n", s, s, len(s), len([]rune(s)))
	r := []rune(s)
	fmt.Printf("[%s] {%#v} len: %d\n", string(r), r, len(r))
	b := []byte(s)
	fmt.Printf("[%s] {%#v} len: %d\n", string(b), b, len(b))
	for index, runeValue := range s {
		fmt.Printf("%#U starts at byte position %d\n", runeValue, index)
	}
}
