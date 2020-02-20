package main

import (
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"reflect"
	"runtime"
	"strings"
	"time"

	"learn.oauth.client/model"
)

var config = struct {
	appID            string
	appSecret        string
	homeURL          string
	authURL          string
	authCodeCallback string
	tokenEndpoint    string
	logoutURL        string
	afterLogoutURL   string
	servicesEndpoint string
}{
	appID:            "billingApp",
	appSecret:        "ff6952f2-7214-468f-ad1d-e45947b75b57",
	homeURL:          "http://localhost:8000",
	authURL:          "http://localhost:8080/auth/realms/learningApp/protocol/openid-connect/auth",
	authCodeCallback: "http://localhost:8000/authCodeRedirect",
	logoutURL:        "http://localhost:8080/auth/realms/learningApp/protocol/openid-connect/logout",
	tokenEndpoint:    "http://localhost:8080/auth/realms/learningApp/protocol/openid-connect/token",
	servicesEndpoint: "http://localhost:8001/billing/v1/services",
	afterLogoutURL:   "http://localhost:8000",
}

var t = template.Must(template.ParseFiles("template/index.html"))
var tServices = template.Must(template.ParseFiles("template/index.html", "template/services.html"))

type AppVar struct {
	AuthCode     string
	SessionState string
	AccessToken  string
	RefreshToken string
	Scope        string
	Services     []string
}

var appVar = AppVar{}

func home(w http.ResponseWriter, r *http.Request) {
	t.Execute(w, appVar)
}

func login(w http.ResponseWriter, r *http.Request) {
	req, err := http.NewRequest("GET", config.authURL, nil)
	if err != nil {
		log.Print(err)
		return
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

func exchangeToken(w http.ResponseWriter, r *http.Request) {
	qs := url.Values{}
	qs.Add("client_id", config.appID)
	qs.Add("grant_type", "authorization_code")
	qs.Add("code", appVar.AuthCode)
	qs.Add("redirect_uri", config.authCodeCallback)

	req, err := http.NewRequest("POST", config.tokenEndpoint, strings.NewReader(qs.Encode()))
	if err != nil {
		log.Print(err)
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.SetBasicAuth(config.appID, config.appSecret)

	c := http.Client{}
	res, err := c.Do(req)
	if err != nil {
		log.Print("couldn't get access token", err)
		return
	}

	byteBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Print(err)
		return
	}
	defer res.Body.Close()

	accessTokenResponse := model.AccessTokenResponse{}
	json.Unmarshal(byteBody, &accessTokenResponse)

	appVar.AccessToken = accessTokenResponse.AccessToken
	appVar.RefreshToken = accessTokenResponse.RefreshToken
	appVar.Scope = accessTokenResponse.Scope

	log.Println("token", appVar.AccessToken)

	t.Execute(w, appVar)
}

func enabledLog(handler func(w http.ResponseWriter, r *http.Request)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		handlerName := runtime.FuncForPC(reflect.ValueOf(handler).Pointer()).Name()
		log.SetPrefix(handlerName + " ")
		log.Printf("--> %s\n", handlerName)
		log.Printf("request: %+v\n", r.RequestURI)
		log.Printf("response: %+v\n", w)
		handler(w, r)
		log.Printf("<-- %s\n\n", handlerName)
	}
}

func main() {
	http.HandleFunc("/", enabledLog(home))
	http.HandleFunc("/login", enabledLog(login))
	http.HandleFunc("/logout", enabledLog(logout))
	http.HandleFunc("/exchangeToken", enabledLog(exchangeToken))
	http.HandleFunc("/services", enabledLog(services))
	http.HandleFunc("/authCodeRedirect", enabledLog(authCodeRedirect))
	http.ListenAndServe(":8000", nil)
}

func services(w http.ResponseWriter, r *http.Request) {
	req, err := http.NewRequest("GET", config.servicesEndpoint, nil)
	if err != nil {
		log.Print(err)
		tServices.Execute(w, appVar)
		return
	}

	req.Header.Add("Authorization", "Bearer "+appVar.AccessToken)

	ctx, cancelFunc := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancelFunc()

	c := http.Client{}
	res, err := c.Do(req.WithContext(ctx))
	if err != nil {
		log.Println(err)
		tServices.Execute(w, appVar)
		return
	}

	byteBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Print(err)
		tServices.Execute(w, appVar)
		return
	}
	defer res.Body.Close()

	billingResponse := &model.BillingResponse{}
	err = json.Unmarshal(byteBody, billingResponse)
	if err != nil {
		log.Print(err)
		tServices.Execute(w, appVar)
		return
	}

	appVar.Services = billingResponse.Services

	tServices.Execute(w, appVar)
}
