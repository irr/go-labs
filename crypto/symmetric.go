package main

import (
    "bytes"
    "fmt"    
    "github.com/cryptobox/gocryptobox/secretbox"
)

func main() {

    message := []byte("This is a test message for secretbox.")

    fmt.Printf("MSG: %#+v\n", string(message))

    key, ok := secretbox.GenerateKey()

    if !ok {
        panic("Failed to generate key!")
    }

    secret, ok := secretbox.Seal(message, key)

    if !ok {
        panic("Failed to seal message!")
    }

    decrypted, ok := secretbox.Open(secret, key)

    if !ok {
        panic("Failed to open message!")
    }

    fmt.Printf("DECRYPTED: %#+v\n", string(decrypted))

    if bytes.Equal(message, decrypted) {
        fmt.Println("secretbox: OK.")
    } else {
        panic("secretbox failed!")
    }
}
