// curl -v -X POST -d "{\"test\": \"that\"}" http://localhost:8082/test
package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type test_struct struct {
	Test string
}

func test(rw http.ResponseWriter, req *http.Request) {
	decoder := json.NewDecoder(req.Body)
	var t test_struct
	err := decoder.Decode(&t)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("%#v\n", t)
}

func main() {
	http.HandleFunc("/test", test)
	log.Fatal(http.ListenAndServe(":8082", nil))
}
