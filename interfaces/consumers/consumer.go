package consumers

import (
	"fmt"
	"sample/engines"
	"time"
)

// Consumer ...
type Consumer struct {
	T string
	E engines.Engine
	F func()
}

// Run ...
func (c Consumer) Run() error {
	fmt.Printf("starting topic: %s\n", c.T)
	for {
		msg, err := c.E.GetMessage()
		if err != nil {
			return err
		}
		err = c.E.Process(&msg)
		if err != nil {
			return err
		}
		c.E.Validate(&msg)

		c.F()

		// Restart
		fmt.Println("waiting...")
		time.Sleep(2 * time.Second)
	}
}
