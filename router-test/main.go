package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"

	"github.com/julienschmidt/httprouter"
)

func other(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id := ps.ByName("id")
	nid := ps.ByName("nid")
	fmt.Fprintf(w, "id: %v and nid: %v\n", id, nid)
}

func headers(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id := ps.ByName("id")
	for name, headers := range r.Header {
		for _, h := range headers {
			fmt.Fprintf(w, "%v: %v\n", name, h)
		}
	}
	fmt.Fprintf(w, "id: %v\n", id)
}

func main() {
	mux := httprouter.New()
	mux.GET("/headers/:id", headers)
	mux.GET("/other/:id/endpoint/:nid", other)

	ts := httptest.NewServer(mux)
	defer ts.Close()

	url := fmt.Sprintf("%s/other/7b2a774e-6d59-11ea-90c6-1f5c859414bc/endpoint/100?q=mytest", ts.URL)

	fmt.Println(url)

	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}

	greeting, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%s", greeting)
}
