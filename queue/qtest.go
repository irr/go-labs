package main

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/enriquebris/goconcurrentqueue"
)

func worker(worknumber int, queue *goconcurrentqueue.FIFO) {
	for true {
		value, err := queue.DequeueOrWaitForNextElement()
		if err != nil {
			log.Panicf("worker error: %v", err)
		}
		fmt.Printf("{worker:%v} got value: %v\n", worknumber, value)
	}
}

func main() {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	queue := goconcurrentqueue.NewFIFO()

	workers := 4
	for i := 0; i < workers; i++ {
		go worker(i, queue)
	}

	wait := 2
	for true {
		n := r.Intn(10) + 1
		fmt.Printf("enqueing %v values...\n", n)
		for i := 0; i < n; i++ {
			v := r.Float64()
			fmt.Printf("enqueing {n:%v}: %v\n", n, v)
			queue.Enqueue(v)
		}
		fmt.Printf("waiting %v secs...\n", wait)
		time.Sleep(time.Duration(wait) * time.Second)
	}
}
