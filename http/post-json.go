package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

func main() {
	client := &http.Client{}
	data := strings.NewReader("{\"test\": \"that\"}")
	req, err := http.NewRequest("POST", "http://localhost:8082/test", data)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Add("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	fmt.Printf("%#v\n", resp)
}
