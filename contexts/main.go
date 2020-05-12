package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"
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
	if n < 1 {
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		time.Sleep(time.Duration(n) * time.Second)
		fmt.Fprintf(w, "%d seconds elapsed", n)
	}
}

func server() {
	http.HandleFunc("/", WaitServer)
	log.Printf("  server: server listening on :%d...\n", PORT)
	http.ListenAndServe(fmt.Sprintf(":%d", PORT), nil)
}

func track(start time.Time, name string) {
	elapsed := time.Since(start)
	log.Printf("  track: %s took %s", name, elapsed)
}

func request(ctx context.Context, wg *sync.WaitGroup, c chan int, l string, t int) {
	start := time.Now()
	defer func() {
		track(start, "request")
		wg.Done()
		log.Printf(" request: exited.\n")
	}()
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
	}
	err = ctx.Err()
	deadline, ok := ctx.Deadline()
	log.Printf(" request: %s got %v [%v,%v:%v]\n", l, resp.Status(), deadline, ok, err)
	if err == nil {
		c <- 1
	}
}

func response(ctx context.Context, wg *sync.WaitGroup, cancel context.CancelFunc, c1 chan int, c2 chan int) {
	start := time.Now()
	defer func() {
		track(start, "response")
		wg.Done()
		log.Printf("response: exited.\n")
	}()
	c := 0
	for c < 2 {
		select {
		case r1 := <-c1:
			log.Printf("response: r1 got %v\n", r1)
			c++
		case r2 := <-c2:
			log.Printf("response: r2 got %v\n", r2)
			c++
		case <-ctx.Done():
			log.Printf("response: ctx.Done() called!\n")
			return
		}
		log.Printf("waiting...\n")
		if c == 2 {
			log.Printf("response: all done!\n")
		}
	}
}

func main() {
	timeout := 4

	log.Printf("    main: started.\n")

	go server()

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)

	wg := sync.WaitGroup{}

	c1 := make(chan int)
	wg.Add(1)
	go request(ctx, &wg, c1, "c1", 2)

	c2 := make(chan int)
	wg.Add(1)
	go request(ctx, &wg, c2, "c2", 7)

	wg.Add(1)
	go response(ctx, &wg, cancel, c1, c2)

	wg.Wait()
	log.Printf("    main: exited.\n")
}
