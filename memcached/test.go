package main

// http://code.google.com/p/memcached/wiki/BinaryProtocolRevamped

import (
    "os"
    "fmt"
    "github.com/dustin/gomemcached/client"
)

func main() {
    conn, err := memcached.Connect("tcp", "localhost:11211")
    if err != nil {
        os.Exit(1)
    }
    r, err := conn.Set(0, "foo", 0, 0, []byte("alessandra.cs@uol.com.br"))
    if err != nil {
        os.Exit(1)
    }
    r, err = conn.Get(0, "foo")
    if err != nil {
        os.Exit(1)
    }
    conn.Close()
    fmt.Printf("foo: %+v\n", r)    
}