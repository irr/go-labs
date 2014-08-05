package main

/*
#cgo CFLAGS: -I/opt/lisp/newlisp
#cgo LDFLAGS: -lnewlisp
#include <stdlib.h>
char* newlispEvalStr(char *cmd);
*/
import "C"

import (
	"fmt"
	"unsafe"
)

const lisp = `[cmd]
(load "/usr/share/newlisp/modules/zlib.lsp")
(bayes-train (parse (lower-case (zlib:gz-read-file "Doyle.txt.gz")) "[^a-z]+" 0) '() '() 'DDB)
(bayes-train '() (parse (lower-case (zlib:gz-read-file "Dowson.txt.gz")) "[^a-z]+" 0) '() 'DDB)
(bayes-train '() '() (parse (lower-case (zlib:gz-read-file "Beowulf.txt.gz")) "[^a-z]+" 0) 'DDB)
(println "Doyle test  : " 
    (bayes-query (parse "adventures of sherlock holmes") 'DDB true)
    (bayes-query (parse "adventures of sherlock holmes") 'DDB))
(println "Downson test: "
    (bayes-query (parse "comedy of masks") 'DDB true)
    (bayes-query (parse "comedy of masks") 'DDB))
(println "Beowulf test: " 
    (bayes-query (parse "hrothgar and beowulf") 'DDB true)
    (bayes-query (parse "hrothgar and beowulf") 'DDB))
[/cmd]`

func main() {
	lib := C.CString("newlispEvalStr")
	defer C.free(unsafe.Pointer(lib))

	code := C.CString(lisp)
	defer C.free(unsafe.Pointer(code))

	res := C.GoString(C.newlispEvalStr(code))

	fmt.Printf("Received %+v characters\n", len(res))
	fmt.Printf("Output:\n%+v\n", res)
}
