package main

import (
	"context"
	"fmt"
	"time"

	"golang.org/x/sync/errgroup"
)

func main() {
	g, ctx := errgroup.WithContext(context.Background())

	// Task 1 — loops and respects ctx
	g.Go(func() error {
		for i := 0; i < 2; i++ {
			select {
			case <-ctx.Done():
				fmt.Println("task 1: cancelled")
				return ctx.Err()
			default:
				fmt.Println("task 1: working...", i)
				time.Sleep(300 * time.Millisecond)
			}
		}
		fmt.Println("task 1 done")
		return nil
	})

	// Task 2 — also respects ctx
	g.Go(func() error {
		for i := 0; i < 10; i++ {
			select {
			case <-ctx.Done():
				fmt.Println("task 2: cancelled")
				return ctx.Err()
			default:
				fmt.Println("task 2: working...", i)
				time.Sleep(500 * time.Millisecond)
			}
		}
		fmt.Println("task 2 done")
		return nil
	})

	// Task 3 — fails early
	g.Go(func() error {
		time.Sleep(700 * time.Millisecond)
		fmt.Println("task 3: failing")
		return fmt.Errorf("task 3 failed")
	})

	if err := g.Wait(); err != nil {
		fmt.Println("Group error:", err)
	}
}
