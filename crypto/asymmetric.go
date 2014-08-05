package main

import (
    "bytes"
    "fmt"    
    "github.com/cryptobox/gocryptobox/box"
)

func main() {

    message := []byte("This is a test message for box.")

    fmt.Printf("MSG: %#+v\n", string(message))

    priv, pub, ok := box.GenerateKey()

    if !ok {
        panic("Failed to generate keys!")
    }

    locked, ok := box.Seal(message, pub)

    if !ok {
        panic("Failed to seal a box!")
    }

    fmt.Printf("%#+v\n", locked)

    decrypted, ok := box.Open(locked, priv)

    if !ok {
        panic("Failed to open a sealed box!")  
    }

    fmt.Printf("MSG: %#+v\n", string(message))

    if bytes.Equal(message, decrypted) {
        fmt.Println("box: OK.")
    } else {
        panic("box failed!")
    }
}
