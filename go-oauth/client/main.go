package main

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
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

	"github.com/google/uuid"
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
	afterLogoutURL:   "http://localhost:8000/home",
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
	State        map[string]struct{}
}

func newAppVar() AppVar {
	return AppVar{State: make(map[string]struct{})}
}

var appVar = newAppVar()

func home(w http.ResponseWriter, r *http.Request) {
	t.Execute(w, appVar)
}

var codeVerifier = "code-challenge43128unreserved-._~nge43128dX"

func makeCodeChallenge(plain string) string {
	h := sha256.Sum256([]byte(plain))
	hs := base64.RawURLEncoding.EncodeToString(h[:])
	return hs
}

func login(w http.ResponseWriter, r *http.Request) {
	req, err := http.NewRequest("GET", config.authURL, nil)
	if err != nil {
		log.Print(err)
		return
	}

	state := uuid.New().String()
	appVar.State[state] = struct{}{}

	qs := url.Values{}
	qs.Add("client_id", config.appID)
	qs.Add("response_type", "code")
	qs.Add("state", state)
	qs.Add("scope", "openid billingService")
	qs.Add("redirect_uri", config.authCodeCallback)
	codeChallenge := makeCodeChallenge(codeVerifier)
	qs.Add("code_challenge", codeChallenge)
	qs.Add("code_challenge_method", "S256")

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
	appVar = newAppVar()

	http.Redirect(w, r, logoutURL.String(), http.StatusFound)
}

func authCodeRedirect(w http.ResponseWriter, r *http.Request) {
	appVar.AuthCode = r.URL.Query().Get("code")

	callBackState := r.URL.Query().Get("state")
	if _, ok := appVar.State[callBackState]; !ok {
		fmt.Fprintf(w, "Error")
		return
	}
	delete(appVar.State, callBackState)

	appVar.SessionState = r.URL.Query().Get("session_state")
	fmt.Printf("Request queries : %+v\n", appVar)

	r.URL.RawQuery = ""
	exchangeToken()

	t.Execute(w, appVar)
}

func exchangeToken() {
	qs := url.Values{}
	qs.Add("client_id", config.appID)
	qs.Add("grant_type", "authorization_code")
	qs.Add("code", appVar.AuthCode)
	qs.Add("redirect_uri", config.authCodeCallback)
	qs.Add("code_verifier", codeVerifier)

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
	http.HandleFunc("/home", enabledLog(home))
	http.HandleFunc("/login", enabledLog(login))
	http.HandleFunc("/logout", enabledLog(logout))
	http.HandleFunc("/services", enabledLog(services))
	http.HandleFunc("/refreshToken", enabledLog(refreshToken))
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

	if res.StatusCode != 200 {
		log.Println(string(byteBody))
		appVar.Services = []string{string(byteBody)}
		tServices.Execute(w, appVar)
		return
	}

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

func refreshToken(w http.ResponseWriter, r *http.Request) {

	qs := url.Values{}
	qs.Add("grant_type", "refresh_token")
	qs.Add("refresh_token", appVar.RefreshToken)

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
		tServices.Execute(w, appVar)
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

	tServices.Execute(w, appVar)
}
