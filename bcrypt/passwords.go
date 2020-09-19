package main

// go get golang.org/x/crypto/bcrypt

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

// HashPassword ...
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	return string(bytes), err
}

// CheckPasswordHash ...
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func main() {
	password := "my secret"
	hash, _ := HashPassword(password) // ignore error for the sake of simplicity

	fmt.Println("Password:", password)
	fmt.Println("Hash:    ", hash)

	match := CheckPasswordHash(password, hash)
	fmt.Println("Match:   ", match)

	rubyhash := "$2a$10$EIJlFSK7ZKWCWvTGKrmd/uNoy/ZVVPrFOaZsQuUlCFyy3yxaJR6Ra"
	match = CheckPasswordHash(password, rubyhash)
	fmt.Println("RubyMatch:", match)
}
