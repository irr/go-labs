package main

import (
	"sample/consumers"
	"sample/engines"
)

func main() {
	consumer := &consumers.Consumer{
		E: &engines.MyEngine{
			Name: "irrlab",
		},
	}

	consumer.Run()
}
