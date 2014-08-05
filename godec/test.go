package main

import (
	"code.google.com/p/godec/dec"
	"fmt"
	"strconv"
)

func main() {
	const x = "3537115888337719"
	const y = "1125899906842624"

	x1, y1 := new(dec.Dec), new(dec.Dec)
	x1.SetString(x)
	y1.SetString(y)

	total1 := new(dec.Dec).QuoExact(x1, y1)
	fmt.Printf("GoDec = %#v\n", total1)

	x2, _ := strconv.ParseFloat(x, 64)
	y2, _ := strconv.ParseFloat(y, 64)
	total2 := x2 / y2
	fmt.Printf("Float = %#v\n", total2)
}
