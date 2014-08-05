package main

/*
$ go install "example/newmath"
compiles and copies binary (.a) to ~/go/golib/pkg/<arch>...

$ cd $GOPATH/src/example/newmath
$ go install

.../pkg/
    linux_amd64/
        example/
            newmath.a  # package object
.../src/
    example/
        newmath/
            sqrt.go    # package source
*/

import (
    "fmt"
    "example/newmath"
)

func main() {
    fmt.Printf("newmath's sqrt92): %v\n", newmath.Sqrt(2))
}
