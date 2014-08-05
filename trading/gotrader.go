package main

import (
    "fmt"
    "os"
)

func main() {
    fmt.Printf("starting gotrader...\n")

    values, err := GetSymbols("UOLL4")

    if err != nil {
        fmt.Printf("GetSymbols error: %#v\n", err)
        os.Exit(1)
    }

    p, n := 10, len(values)
    for i := 0; i < n; i++ {
        fmt.Printf("%05d: %v\n", i+1, values[i])
        if i == (p - 1) {
            fmt.Printf("...\n%05d: %v\n", n, values[n-1])
            break
        }
    }

    fmt.Printf("exiting gotrader...\n")
}
