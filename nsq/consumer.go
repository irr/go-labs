package main

import (
  "log"
  "github.com/nsqio/go-nsq"
  "sync"
  "time"
)

func read(wg *sync.WaitGroup) {
    config := nsq.NewConfig()
    q, _ := nsq.NewConsumer("write_test", "ch", config)
    q.AddHandler(nsq.HandlerFunc(func(message *nsq.Message) error {
      log.Printf("Got: %+v => message: %v", message, string(message.Body[:]))
      return nil
    }))
    err := q.ConnectToNSQLookupd("127.0.0.1:4161")
    if err != nil {
      wg.Done()
      log.Panic("Could not connect")
    }
    for {
        log.Printf("Stats: (%+v)=>%+v", &q, q.Stats())
        time.Sleep(1000 * time.Millisecond)
    }
}
func main() {
    wg := &sync.WaitGroup{}
    wg.Add(1)
    go read(wg)
    wg.Wait()
}
