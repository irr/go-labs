package main

// This is an example of a resilient service worker program written in Go.
//
// This program will run a worker every 5 seconds and exit when SIGINT or SIGTERM
// is received, while ensuring any ongoing work is finished before exiting.
//
// Unexpected panics are also handled: program won't crash if the worker panics.

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"runtime/debug"
	"sync"
	"syscall"
	"time"
)

var (
	sigChan    chan os.Signal
	workerChan chan bool
	waitGroup  sync.WaitGroup
	errLogger  *log.Logger
)

const sleepDuration time.Duration = time.Duration(1) * time.Second

func runWorker() {
	defer recoverWorker()
	waitGroup.Add(1)

	// Fake long process
	time.Sleep(sleepDuration)

	waitGroup.Done()
	workerChan <- true
}

func recoverWorker() {
	if err := recover(); err != nil {
		// Handle unexpected panic
		errLogger.Println(err)
		errLogger.Print(string(debug.Stack()))

		// Finish worker execution anyway
		waitGroup.Done()

		// Return false if service should stop on panic
		// Service will continue otherwise
		workerChan <- true
	}
}

func runTimer() {
	for {
		workerChan <- true
		if !<-workerChan { // Exit if worker returned false
			sigChan <- syscall.SIGTERM
			return
		}
		fmt.Println("OK")
		time.Sleep(sleepDuration)
	}
}

func listen() {
	for {
		select {
		case <-sigChan:
			// Wait for worker to finish before exit
			waitGroup.Wait()
			return
		case <-workerChan:
			go runWorker()
		}
	}
}

func setup() {
	sigChan = make(chan os.Signal)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	workerChan = make(chan bool)

	errLogger = log.New(os.Stderr, "", log.LstdFlags)
}

func main() {
	setup()

	go runTimer()
	fmt.Println("Service running...")
	listen()
	fmt.Println("Service stopped.")
}