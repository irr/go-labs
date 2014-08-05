// curl -v -X POST -d "{\"test\": \"that\"}" http://localhost:8082/test
package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

func test(rw http.ResponseWriter, req *http.Request) {
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(string(body))
	var m map[string]interface{}
	err = json.Unmarshal(body, &m)
	if err != nil {
		log.Println(err)
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}
	log.Printf("%#v\n", m)
}

func main() {
	http.HandleFunc("/test", test)
	log.Fatal(http.ListenAndServe(":8082", nil))
}
