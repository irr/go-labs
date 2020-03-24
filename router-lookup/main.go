package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"

	"github.com/julienschmidt/httprouter"
)

func httprouterHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id := ps.ByName("id")
	nid := ps.ByName("nid")
	fmt.Fprintf(w, "id: %v and nid: %v\n", id, nid)
}

func main() {
	mux := httprouter.New()
	mux.GET("/other/:id/endpoint/:nid", httprouterHandler)
	mux.GET("/other/:id/endpoint/:nid/", httprouterHandler)

	uri := "/other/7b2a774e-6d59-11ea-90c6-1f5c859414bc/endpoint/100?q=mytest"

	u, err := url.Parse(uri)
	if err != nil {
		panic(err)
	}

	fmt.Printf(" u:%+v\n", u.Path)

	w := httptest.NewRecorder()
	h, p, ok := mux.Lookup("GET", u.Path)

	if h != nil {
		fmt.Printf(" h: %+v\n p: %+v\nok: %+v\n", h, p, ok)

		h(w, nil, p)

		fmt.Printf("%+v\n", w.Body.String())
	}
}
