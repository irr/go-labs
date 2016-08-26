package main

// curl -d "message from nsq1" http://127.0.0.1:4151/pub?topic=write_test

import (
    "fmt"
    "log"
    "github.com/nsqio/go-nsq"
    "sync"
)

func pub(wg *sync.WaitGroup, port int) {
    defer wg.Done()
    config := nsq.NewConfig()
    w, _ := nsq.NewProducer(fmt.Sprintf("127.0.0.1:%d", port), config)
    err := w.Publish("write_test", []byte(fmt.Sprintf("message from queue {127.0.0.1:%d}", port)))
    if err != nil {
      log.Panic("Could not connect")
    }
    w.Stop()
}

func main() {
    wg := &sync.WaitGroup{}
    ports := []int{4150, 4250}
    for _, port := range ports {
        wg.Add(1)
        go pub(wg, port)
    }
    wg.Wait()
}
