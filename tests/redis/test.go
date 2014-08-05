package main

// https://github.com/simonz05/godis
import (
    "godis"
    "fmt"
)

func main() {
    // new client on default port 6379, select db 0 and use no password
    c := godis.New("tcp:127.0.0.1:6379", 0, "")

    // set the key "foo" to "Hello Redis"
    c.Set("foo", "Hello Redis")

    // retrieve the value of "foo". Returns an Elem obj
    elem, _ := c.Get("foo")

    // convert the obj to a string and print it
    fmt.Println("foo: ", elem.String())
}
