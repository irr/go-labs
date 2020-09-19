package main

import (
	"fmt"
	"sample/consumers"
	"sample/engines"
	"time"
)

func main() {
	consumer := &consumers.Consumer{
		T: "topic",
		E: &engines.MyEngine{
			Name: "irrlab",
		},
		F: func() { fmt.Println("F") },
	}

	go consumer.Run()

	consumer2 := &consumers.Consumer{
		T: "topic2",
		E: &engines.MyEngine2{
			Name:  "irrlab",
			Other: "other",
		},
		F: func() { fmt.Println("F2") },
	}

	go consumer2.Run()

	time.Sleep(100 * time.Second)
}
