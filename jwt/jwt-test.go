package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

// SM ...
const (
	SM     = "RS512"
	SECRET = "anypassword"
)

// CustomClaims ...
type CustomClaims struct {
	Foo string `json:"foo"`
	jwt.StandardClaims
}

// Generate keys using: https://www.csfieldguide.org.nz/en/interactives/rsa-key-generator/
// Keysize: 1024 and PKCS #8 (base64)
func generateRSAToken(prv string) (string, error) {
	keyData, err := ioutil.ReadFile(prv)
	if err != nil {
		return "", err
	}

	key, err := jwt.ParseRSAPrivateKeyFromPEM(keyData)
	if err != nil {
		return "", err
	}

	token := jwt.New(jwt.GetSigningMethod(SM))
	token.Claims = CustomClaims{
		"bar",
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
			Issuer:    "test",
		},
	}

	return token.SignedString(key)
}

func verifyRSAToken(token, pub string) (bool, error) {
	keyData, err := ioutil.ReadFile(pub)
	if err != nil {
		return false, err
	}

	key, err := jwt.ParseRSAPublicKeyFromPEM(keyData)
	if err != nil {
		return false, err
	}

	parts := strings.Split(token, ".")
	method := jwt.GetSigningMethod(SM)
	err = method.Verify(strings.Join(parts[0:2], "."), parts[2], key)
	if err != nil {
		return false, nil
	}

	return true, nil
}

func generateHASHToken(pwd string) (string, error) {
	claims := CustomClaims{
		"bar",
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
			Issuer:    "test",
		},
	}
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return at.SignedString([]byte(pwd))
}

func verifyHASHToken(token, pwd string) (bool, error) {
	_, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(pwd), nil
	})
	return err == nil, err
}

func main() {
	pub := "./pub.key"
	prv := "./prv.key"

	token, err := generateRSAToken(prv)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("RSA token: %s\n", token)

	valid, err := verifyRSAToken(token, pub)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("RSA token validation: %v (%v)\n", valid, err)

	token, err = generateHASHToken(SECRET)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("HASH token: %s\n", token)

	valid, err = verifyHASHToken(token, SECRET)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("RSA token validation: %v (%v)\n", valid, err)
}
