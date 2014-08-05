package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
)

func main() {
	urls := []string{
		"http://www.uol.com.br",
		"http://www.google.com",
		"http://www.yahoo.com",
		"http://br.baidu.com",
		"http://www.amazon.com"}
	var wg sync.WaitGroup
	defer wg.Wait()
	barrier := make(chan int, 3)
	for _, v := range urls {
		url := v
		barrier <- 1
		wg.Add(1)
		go func() {
			defer func() { wg.Done(); <-barrier }()
			res, err := http.Get(url)
			if err != nil {
				log.Fatal(err)
			}
			body, err := ioutil.ReadAll(res.Body)
			defer res.Body.Close()
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("url: %-25s %6d bytes\n", url, len(body))
		}()
	}
}
