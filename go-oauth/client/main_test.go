package main

import (
	"log"
	"testing"
)

func TestMakeCodeChallenge(t *testing.T) {
	codeVerifier := "code-challenge43128unreserved-._~nge43128dX"
	challenge := makeCodeChallenge(codeVerifier)
	log.Println(challenge)
}
