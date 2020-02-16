package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"
)

var config = struct {
	appID            string
	homeURL          string
	authURL          string
	authCodeCallback string
	logoutURL        string
	afterLogoutURL   string
}{
	appID:            "billingApp",
	homeURL:          "http://localhost:8000",
	authURL:          "http://localhost:8080/auth/realms/learningApp/protocol/openid-connect/auth",
	authCodeCallback: "http://localhost:8000/authCodeRedirect",
	logoutURL:        "http://localhost:8080/auth/realms/learningApp/protocol/openid-connect/logout",
	afterLogoutURL:   "http://localhost:8000",
}

var t = template.Must(template.ParseFiles("template/index.html"))

type AppVar struct {
	AuthCode     string
	SessionState string
}

var appVar = AppVar{}

func home(w http.ResponseWriter, r *http.Request) {
	t.Execute(w, appVar)
}

func login(w http.ResponseWriter, r *http.Request) {
	req, err := http.NewRequest("GET", config.authURL, nil)
	if err != nil {
		log.Print(err)
	}

	qs := url.Values{}
	qs.Add("client_id", config.appID)
	qs.Add("response_type", "code")
	qs.Add("state", "123")
	qs.Add("redirect_uri", config.authCodeCallback)

	req.URL.RawQuery = qs.Encode()
	http.Redirect(w, r, req.URL.String(), http.StatusFound)
}

func logout(w http.ResponseWriter, r *http.Request) {
	qs := url.Values{}
	qs.Add("redirect_uri", config.afterLogoutURL)

	logoutURL, err := url.Parse(config.logoutURL)
	if err != nil {
		log.Println(err)
	}
	logoutURL.RawQuery = qs.Encode()
	appVar = AppVar{}

	http.Redirect(w, r, logoutURL.String(), http.StatusFound)
}

func authCodeRedirect(w http.ResponseWriter, r *http.Request) {
	appVar.AuthCode = r.URL.Query().Get("code")
	appVar.SessionState = r.URL.Query().Get("session_state")
	fmt.Printf("Request queries : %+v\n", appVar)

	r.URL.RawQuery = ""
	http.Redirect(w, r, config.homeURL, http.StatusFound)
	t.Execute(w, nil)
}

func main() {
	http.HandleFunc("/", home)
	http.HandleFunc("/login", login)
	http.HandleFunc("/logout", logout)
	http.HandleFunc("/authCodeRedirect", authCodeRedirect)
	http.ListenAndServe(":8000", nil)
}
