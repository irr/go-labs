package main

import (
	"fmt"

	"github.com/sbstjn/allot"
)

func main() {
	cmd := allot.New("/manager/<id:integer>/entity/<name:string>/class/(private|public)*<dummy:string>*")
	match, err := cmd.Match("/manager/1972/entity/irrlab/class/private/?q=1")

	if err != nil {
		fmt.Println("Request did not match command.", err)
	} else {
		fmt.Printf("%+v\n", match)

		id, _ := match.Integer("id")
		name, _ := match.String("name")
		class, _ := match.Match(2)
		dummy, _ := match.Match(3)

		fmt.Printf("Id \"%d\" on \"%s\" at \"%s\" (%s)\n", id, name, class, dummy)
	}
}
