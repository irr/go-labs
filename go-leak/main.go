package main

import (
	"fmt"
)

func someFunc() error {
	return nil
}

// Run ...
func Run() {
	c := make(chan error)
	go func() {
		for err := range c {
			if err != nil {
				panic(err)
			} else {
				fmt.Printf("OK [%+v]\n", err)
			}
		}
	}()
	c <- someFunc()
	close(c)
}

func main() {
	Run()
}
