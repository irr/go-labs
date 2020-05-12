package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/go-resty/resty/v2"
)

// PORT ...
var PORT = 7070

// WaitServer ...
func WaitServer(w http.ResponseWriter, r *http.Request) {
	secs := r.URL.Query().Get("secs")
	n, err := strconv.Atoi(secs)
	if err != nil {
		n = 1
	}
	time.Sleep(time.Duration(n) * time.Second)
	fmt.Fprintf(w, "%d seconds elapsed", n)
}

func server() {
	http.HandleFunc("/", WaitServer)
	log.Printf("  server: server listening on :%d...\n", PORT)
	http.ListenAndServe(fmt.Sprintf(":%d", PORT), nil)
}

func request(ctx context.Context, c chan int, l string, t int) {
	go func() {
		client := resty.New()
		log.Printf(" request: %s sent with t:%d\n", l, t)
		resp, err := client.R().
			SetQueryParams(map[string]string{
				"secs": fmt.Sprintf("%d", t),
			}).
			SetContext(ctx).
			Get(fmt.Sprintf("http://localhost:%d", PORT))

		if err != nil {
			log.Printf(" request: %s error [%v]\n", l, err)
			c <- 0
			return
		}
		err = ctx.Err()
		deadline, ok := ctx.Deadline()
		log.Printf(" request: %s got %v [%v,%v:%v]\n", l, resp.Status(), deadline, ok, err)
		if err == nil {
			c <- 1
		}
	}()
}

func response(ctx context.Context, cancel context.CancelFunc, c1 chan int, c2 chan int) {
	defer func() {
		log.Printf("response: exited.\n")
	}()
	for {
		select {
		case r1 := <-c1:
			log.Printf("response: r1 got %v\n", r1)
		case r2 := <-c2:
			log.Printf("response: r2 got %v\n", r2)
		case <-ctx.Done():
			log.Printf("response: ctx.Done() called!\n")
			return
		}
	}
}

func main() {
	max, timeout := 5, 3

	log.Printf("    main: started.\n")

	go server()

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)

	c1 := make(chan int)
	go request(ctx, c1, "c1", 2)

	c2 := make(chan int)
	go request(ctx, c2, "c2", 4)

	go response(ctx, cancel, c1, c2)

	time.Sleep(time.Duration(max) * time.Second)

	log.Printf("    main: exited.\n")
}
