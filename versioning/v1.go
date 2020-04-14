package main

import (
	"encoding/json"
	"net/http"
)

type Profile struct {
	Name    string
	Hobbies []string
}

func main() {
	http.HandleFunc("/", foo)
	http.ListenAndServe(":3001", nil)
}

func foo(w http.ResponseWriter, r *http.Request) {
	profile := Profile{"Version", []string{"v1"}}

	js, err := json.Marshal(profile)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/vnd.t.v1+json")

	w.Write(js)
}
