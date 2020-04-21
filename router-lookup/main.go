package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"runtime/pprof"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/oklog/ulid"
)

func httprouterHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id := ps.ByName("id")
	nid := ps.ByName("nid")
	fmt.Fprintf(w, "id: %v and nid: %v\n", id, nid)
}

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")

func main() {
	flag.Parse()
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	mux := httprouter.New()
	mux.GET("/other/:id/endpoint/:nid", httprouterHandler)
	mux.GET("/other/:id/endpoint/:nid/", httprouterHandler)

	entropy := ulid.Monotonic(rand.New(rand.NewSource(time.Now().UnixNano())), 0)

	for i := 1; i < 10000; i++ {
		uri := fmt.Sprintf("/other/%v/endpoint/100?q=mytest", ulid.MustNew(ulid.Timestamp(time.Now()), entropy))

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
}
