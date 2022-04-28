package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"
)

func server() {
	// Create an HTTP server that listens on port 8000
	http.ListenAndServe(":8000", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		// This prints to STDOUT to show that processing has started
		fmt.Fprint(os.Stdout, "processing request\n")
		// We use `select` to execute a piece of code depending on which
		// channel receives a message first
		select {
		case <-time.After(2 * time.Second):
			// If we receive a message after 2 seconds
			// that means the request has been processed
			// We then write this as the response
			w.Write([]byte("request processed"))
		case <-ctx.Done():
			// If the request gets cancelled, log it
			// to STDERR
			fmt.Fprint(os.Stderr, "request cancelled\n")
		}
	}))
}

func main() {

	go server()

	ctx, cancel := context.WithCancel(context.Background())

	// Make a request, that will call the google homepage
	req, err := http.NewRequest(http.MethodGet, "http://localhost:8000/", nil)

	// Associate the cancellable context we just created to the request
	req = req.WithContext(ctx)

	go func() {
		time.Sleep(1 * time.Second)
		cancel()
	}()

	// Create a new HTTP client and execute the request
	client := &http.Client{}
	res, err := client.Do(req)

	// If the request failed, log to STDOUT
	if err != nil {
		fmt.Println("Request failed:", err)
	} else {
		// Print the statuscode if the request succeeds
		fmt.Println("Response received, status code:", res.StatusCode)
	}

	time.Sleep(5 * time.Second)
}
