package main

// #include <stdio.h>
// #include <stdlib.h>
import "C"
import "unsafe"

import ( "fmt" )

// http://golang.org/cmd/cgo/
// http://golang.org/doc/articles/c_go_cgo.html

func print(s string) {
    cs := C.CString(s)
    defer C.free(unsafe.Pointer(cs))
    C.fputs(cs, (*C.FILE)(C.stdout))
    C.fflush((*C.FILE)(C.stdout))
}

// go build test.go
// valgrind --leak-check=full test

func main() {
	print("test ok! (print)\n")
	fmt.Println("test ok! (fmt)")
}