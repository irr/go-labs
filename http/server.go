package main

import (
	"flag"
	"fmt"
	"log"
	"log/syslog"
	"net/http"
	"os"
)

var l *bool

var logger *syslog.Writer
var err error

func T(exp bool, a interface{}, b interface{}) interface{} {
	if exp {
		return a
	}
	return b
}

func getValue(req *http.Request, n string, d string) string {
	v := req.FormValue(n)
	return (T((v == ""), d, v)).(string)
}

func checkError(err error, w http.ResponseWriter, msg string, status int) bool {
	if err != nil {
		log.Printf("%s: %s\n", msg, err.Error())
		http.Error(w, err.Error(), status)
	}
	return (err != nil)
}

func TestServer(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, "%s\n", getValue(req, "name", "irr"))
}

func main() {
	l = flag.Bool("l", false, "syslog enabled/disabled")
	h := flag.Bool("h", false, "help")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "usage: server [-l][-h]\n")
		flag.PrintDefaults()
	}

	flag.Parse()

	if *h {
		flag.Usage()
		os.Exit(0)
	}

	if *l {
		logger, err = syslog.New(syslog.LOG_INFO, "[test-redigo]")
		if err != nil {
			log.Fatal("syslog: ", err)
		}
	}

	http.HandleFunc("/test", TestServer)
	err := http.ListenAndServe(":3000", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
