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
	log.Printf("server listening on :%d...\n", PORT)
	http.ListenAndServe(fmt.Sprintf(":%d", PORT), nil)
}

func request(c chan int, l string, t int) {
	go func() {
		client := resty.New()
		resp, err := client.R().
			SetQueryParams(map[string]string{
				"secs": fmt.Sprintf("%d", t),
			}).
			EnableTrace().
			Get(fmt.Sprintf("http://localhost:%d", PORT))
		if err != nil {
			log.Printf("%s error [%v]\n", l, err)
			c <- 0
			return
		}
		log.Printf("%s got %v\n", l, resp.Status())
		c <- 1
	}()
}

func response(ctx context.Context, cancel context.CancelFunc, c1 chan int, c2 chan int) {
	for {
		select {
		case r1 := <-c1:
			log.Printf("r1 got %v\n", r1)
			//cancel()
		case r2 := <-c2:
			log.Printf("r2 got %v\n", r2)
		case <-ctx.Done():
			log.Printf("ctx.Done() called!\n")
			return
		}
	}
}

func main() {
	max := 5
	timeout := 2

	go server()

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)

	c1 := make(chan int)
	request(c1, "c1", 1)

	c2 := make(chan int)
	request(c2, "c2", 3)

	go response(ctx, cancel, c1, c2)

	time.Sleep(time.Duration(max) * time.Second)
	log.Println("server exited.")
}
