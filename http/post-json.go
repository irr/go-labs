package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
)

func main() {
	client := &http.Client{}
	v := url.Values{}
	v.Add("test", "One")
	v.Add("test", "Two")
	b, err := json.Marshal(v)
	if err != nil {
		log.Fatal(err)
	}
	data := strings.NewReader(string(b))
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
