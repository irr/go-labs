package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"

	"github.com/gojektech/heimdall/v6/httpclient"
)

func main() {

	target := "https://wttr.in/lisbon?%s"
	client := httpclient.NewClient()

	params := url.Values{}
	params.Add("format", "3")
	params.Add("lang", "en")
	q := params.Encode()

	res, err := client.Get(fmt.Sprintf(target, q), nil)
	if err != nil {
		log.Fatal(err)
	}

	body, err := ioutil.ReadAll(res.Body)
	defer res.Body.Close()

	if res.StatusCode != 200 {
		log.Fatal(res)
	}

	fmt.Printf("%s\n", body)

	req, err := http.NewRequest("GET", "https://wttr.in/lisbon", nil)
	values := req.URL.Query()
	values.Add("format", "3")
	req.URL.RawQuery = values.Encode()

	res, err = client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	body, err = ioutil.ReadAll(res.Body)
	defer res.Body.Close()

	if res.StatusCode != 200 {
		log.Fatal(res)
	}

	fmt.Printf("URL      %+v\n", req.URL)
	fmt.Printf("RawQuery %+v\n", req.URL.RawQuery)
	fmt.Printf("Query    %+v\n", req.URL.Query())

	fmt.Printf("%s\n", body)
}
