package consumers

import (
	"fmt"
	"sample/engines"
	"time"
)

// Consumer ...
type Consumer struct {
	E engines.Engine
}

// Run ...
func (c Consumer) Run() error {
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

		fmt.Println("waiting...")
		time.Sleep(2 * time.Second)
	}
}
