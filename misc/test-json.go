package main

import (
	"encoding/json"
	"fmt"
)

var jsonBlob = []byte(`[
    {"Name": "Platypus", "Order": "Monotremata"},
    {"Name": "Quoll" },
	{"Name": "Nulus", "Order": null, "Test":100},
	{"Name": "Alien", "Order": "Unknown" }
]`)

type Animal struct {
	Name  string
	Order *string
}

func main() {
	var unknown string = "null"
	var animals []Animal
	err := json.Unmarshal(jsonBlob, &animals)
	if err != nil {
		fmt.Println("error:", err)
	} else {
		fmt.Printf("(%d): %+v\n", len(animals), animals)
		for _, v := range animals {
			if v.Order == nil {
				v.Order = &unknown
			}
			fmt.Printf("%s: %+v(@%+v)\n", v.Name, *v.Order, v.Order)
			b, err := json.Marshal(v)
			if err != nil {
				fmt.Println("error:", err)
			} else {
				fmt.Printf("\tjson => %s\n", b)
			}
		}
	}
}
