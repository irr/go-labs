package main

import (
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/gorilla/mux"

	"go.keycloak.svr/procon_data"
)

var home_tpl = template.Must(template.ParseFiles("templates/index.html"))
var systems_tpl = template.Must(template.ParseFiles("templates/index.html", "templates/systems.html"))

var addr = flag.String("addr", "0.0.0.0:8000", "http service address")

var config = struct {
	clientId                       string
	clientSecret                   string
	authUrl                        string
	authCodeCallback               string
	logoutUrl                      string
	afterLogoutRedirect            string
	tokenEndpoint                  string
	validateToken                  string
	appTokenValidationClient       string
	appTokenValidationClientSecret string
}{
	clientId:                       "TestingApp",
	clientSecret:                   "200bb020-4a33-487f-b0b0-d9a23697ffa0",
	authUrl:                        "http://localhost:8080/auth/realms/TestRealm/protocol/openid-connect/auth",
	authCodeCallback:               "http://localhost:8000/authCodeRedirect",
	logoutUrl:                      "http://localhost:8080/auth/realms/TestRealm/protocol/openid-connect/logout",
	afterLogoutRedirect:            "http://localhost:8000",
	tokenEndpoint:                  "http://localhost:8080/auth/realms/TestRealm/protocol/openid-connect/token",
	validateToken:                  "http://localhost:8080/auth/realms/TestRealm/protocol/openid-connect/token/introspect",
	appTokenValidationClient:       "TestingAppValidateTokenClient",
	appTokenValidationClientSecret: "a19c1348-0c24-419a-9dec-804f8e0f817d",
}

type AppVars struct {
	AuthCode     string
	SessionState string

	AccessToken  string
	RefreshToken string
	Scope        string

	Systems []struct {
		Host string `json:"host"`
		Port string `json:"port"`
	} `json:"systems"`
}

var appVars = AppVars{}

func home(w http.ResponseWriter, r *http.Request) {
	home_tpl.Execute(w, appVars)
}

func login(w http.ResponseWriter, r *http.Request) {
	req, err := http.NewRequest("GET", config.authUrl, nil)
	if err != nil {
		fmt.Println(err)
	} else {
		qp := url.Values{}
		qp.Add("state", "noop")
		qp.Add("client_id", config.clientId)
		qp.Add("response_type", "code")
		qp.Add("redirect_uri", config.authCodeCallback)
		req.URL.RawQuery = qp.Encode()
		http.Redirect(w, r, req.URL.String(), http.StatusFound)
	}
}

func logout(w http.ResponseWriter, r *http.Request) {
	qp := url.Values{}
	qp.Add("redirect_uri", config.afterLogoutRedirect)

	logoutUrl, err := url.Parse(config.logoutUrl)
	logoutUrl.RawQuery = qp.Encode()

	appVars = AppVars{}

	if err != nil {
		fmt.Println(("Error parsing logout url"))
	} else {
		http.Redirect(w, r, logoutUrl.String(), http.StatusFound)
	}
}

func authCodeRedirect(w http.ResponseWriter, r *http.Request) {
	appVars.AuthCode = r.URL.Query().Get("code")
	appVars.SessionState = r.URL.Query().Get("session_state")
	r.URL.RawQuery = ""
	fmt.Printf("Req: %+v \n", appVars)
	http.Redirect(w, r, "http://localhost:8000", http.StatusFound)
}

func exchangeToken(w http.ResponseWriter, r *http.Request) {
	form := url.Values{}
	form.Add("grant_type", "authorization_code")
	form.Add("code", appVars.AuthCode)
	form.Add("redirect_uri", config.authCodeCallback)
	form.Add("client_id", config.clientId)

	req, err := http.NewRequest("POST", config.tokenEndpoint, strings.NewReader(form.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.SetBasicAuth(config.clientId, config.clientSecret)
	if err != nil {
		fmt.Println(("Error initializing exchange request"))
	} else {
		c := http.Client{}
		res, err := c.Do(req)
		if err != nil {
			fmt.Println(("Error with doing access token request"))
		} else {
			defer res.Body.Close()
			data, err := ioutil.ReadAll(res.Body)
			if err != nil {
				fmt.Println("Error reading access token response")
			} else {
				at := &procon_data.AccessToken{}
				json.Unmarshal(data, at)
				appVars.AccessToken = at.AccessToken
				appVars.RefreshToken = at.RefreshToken
				appVars.Scope = at.Scope

				home_tpl.Execute(w, appVars)
			}
		}
	}
}

func extractToken(r *http.Request) (string, error) {
	header_token := r.Header.Get("Authorization")
	body_token := r.FormValue("access_token")
	query_token := r.URL.Query().Get("access_token")
	token := ""
	switch {
	case header_token != "":
		split_auth_header := strings.Split(header_token, " ")
		if len(split_auth_header) != 2 {
			return "", fmt.Errorf("Invalid Authorization header format [%s]", header_token)
		}
		token = split_auth_header[1]
	case body_token != "":
		token = body_token
	case query_token != "":
		token = query_token
	default:
		return "", fmt.Errorf("No Access Token")
	}

	if token != "" {
		return token, nil
	} else {
		return "", fmt.Errorf("No Access Token")
	}
}

func validateToken(token string) bool {
	form := url.Values{}
	form.Add("token", token)
	form.Add("token_type_hint", "requesting_party_token")

	req, err := http.NewRequest("POST", config.validateToken, strings.NewReader(form.Encode()))
	if err != nil {
		fmt.Println(err)
	} else {
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		req.SetBasicAuth(config.appTokenValidationClient, config.appTokenValidationClientSecret)

		c := http.Client{}
		res, err := c.Do(req)
		if err != nil {
			fmt.Println(err)
		} else {
			defer res.Body.Close()
			data, err := ioutil.ReadAll(res.Body)
			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Println(string(data))
				introSpect := &procon_data.TokenInstrospect{}
				json.Unmarshal(data, introSpect)
				if err != nil {
					fmt.Println(err)
				} else {
					fmt.Println("INTROSPECT-RESULT: ", introSpect)
					return introSpect.Active
				}
			}
		}
	}
	return false
}

func extractValidateClaims(token string) bool {
	tokenParts := strings.Split(token, ".")
	claim, err := base64.RawURLEncoding.DecodeString(tokenParts[1])
	if err != nil {
		fmt.Println(err)
		return false
	}
	tokenClaim := &procon_data.Tokenclaim{}
	err = json.Unmarshal(claim, tokenClaim)
	if err != nil {
		fmt.Println(err)
		return false
	}

	if !strings.Contains(tokenClaim.Scope, "getSystems") {
		return false
	}

	fmt.Println("TOKENCLAIM-RESULT: ", string(claim))
	fmt.Println("TOKENCLAIM-AUDIENCE: ", tokenClaim.AudAsSlice())

	isValidAudience := false
	for _, v := range tokenClaim.AudAsSlice() {
		if v == "TestingApp" || v == "implicitReactClient" {
			isValidAudience = true
			break
		}
	}

	return isValidAudience
}

func apiResourceSystems(w http.ResponseWriter, r *http.Request) {
	hiveData := []byte(``)

	token, err := extractToken(r)

	if err != nil {
		fmt.Println(err)
		hiveData = []byte(`[{"host":"Error","port":"Invalid Token @Extract"}]`)
	} else {
		fmt.Println("Token:" + token)
		if !validateToken(token) {
			hiveData = []byte(`[{"host":"Error","port":"Invalid Token @Validate"}]`)
		} else {
			if !extractValidateClaims(token) {
				hiveData = []byte(`[{"host":"Error","port":"Invalid Token @ValidateClaims"}]`)
			} else {
				hiveData = []byte(`[
			{
				"host": "localhost",
				"port": "8080"
			},
			{
				"host": "something",
				"port": "7000"
			},
			{
				"host": "somewhere",
				"port": "7000"
			}
		]`)
			}
		}
	}
	w.Header().Add("Content-Type", "application/json")
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Write(hiveData)
}

func getSystems(w http.ResponseWriter, r *http.Request) {
	req, err := http.NewRequest("GET", "http://localhost:8000/api/resource/systems", nil)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(req)
		req.Header.Add("Authorization", "Bearer "+appVars.AccessToken)
		c := http.Client{}
		res, err := c.Do(req)
		if err != nil {
			fmt.Println("Error with req @ getSystems")
		} else {
			data, err := ioutil.ReadAll(res.Body)
			if err != nil {
				fmt.Println("Error reading result @getSystems", err)
			} else {
				defer res.Body.Close()
				hive := &procon_data.Hive{}
				err := json.Unmarshal(data, &hive.Systems)
				if err != nil {
					fmt.Println("Error unpacking data to object @getSystems", err)
				} else {
					appVars.Systems = hive.Systems
					fmt.Println(appVars.Systems)
					systems_tpl.Execute(w, appVars)
					return
				}
			}
		}
	}
	systems_tpl.Execute(w, appVars)
}

func main() {
	fmt.Println("Initializing Server...")

	r := mux.NewRouter()
	r.HandleFunc("/", home)
	r.HandleFunc("/login", login)
	r.HandleFunc("/logout", logout)
	r.HandleFunc("/authCodeRedirect", authCodeRedirect)
	r.HandleFunc("/exchange", exchangeToken)

	r.HandleFunc("/systems", getSystems)
	r.HandleFunc("/api/resource/systems", apiResourceSystems)

	fmt.Println("Server Running on port: 8000")
	http.ListenAndServe(*addr, r)
}
