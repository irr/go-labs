package main

import (
    "fmt"
    "github.com/AllenDang/simhash"
)

func main() {
    s1 := "Golang - mapping an variable length array to a string"
    s2 := "Golang - mapping an variable length array to a struct"
    
    likeness := simhash.GetLikenessValue(s1, s2)

    fmt.Println("Likeness:", likeness)
}
