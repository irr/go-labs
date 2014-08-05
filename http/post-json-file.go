package main

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	inFile, err := os.Open("test.json")
	if err != nil {
		log.Fatal(err)
	}
	defer inFile.Close()
	client := &http.Client{}
	req, err := http.NewRequest("POST", "http://localhost:8082/test", bufio.NewReader(inFile))
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
