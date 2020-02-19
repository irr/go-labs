package main

import (
	"fmt"

	"github.com/Nerzal/gocloak/v4"
)

func main() {

	userName := "Bob"
	firstName := "Bob"
	lastName := "Uncle"
	email := "bob@uncle.com"
	enabled := true

	client := gocloak.NewClient("http://localhost:8080")
	token, err := client.LoginAdmin("admin", "admin", "master")
	if err != nil {
		fmt.Printf("%+v\n", err)
		panic("Something wrong with the credentials or url")
	}
	user := gocloak.User{
		FirstName: &firstName,
		LastName:  &lastName,
		Email:     &email,
		Enabled:   &enabled,
		Username:  &userName,
	}
	res, err := client.CreateUser(token.AccessToken, "AuthTest", user)
	if err != nil {
		panic("Oh no!, failed to create user :(")
	}
	fmt.Printf("%+v\n", res)
}
