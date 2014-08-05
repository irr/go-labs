package main

// go get -v code.google.com/p/go.text/transform

import (
	"bytes"
	"fmt"
	"unicode"
	"unicode/utf8"

	"code.google.com/p/go.text/transform"
	"code.google.com/p/go.text/unicode/norm"
)

var isMn = func(r rune) bool {
	return unicode.Is(unicode.Mn, r) // Mn: nonspacing marks
}

var transliterations = map[rune]string{
	'Æ': "AE", 'Ð': "D", 'Ł': "L", 'Ø': "OE", 'Þ': "Th",
	'ß': "ss", 'æ': "ae", 'ð': "d", 'ł': "l", 'ø': "oe",
	'þ': "th", 'Œ': "OE", 'œ': "oe",
}

func RemoveAccents(b []byte) ([]byte, error) {
	mnBuf := make([]byte, len(b)*125/100)
	t := transform.Chain(norm.NFD, transform.RemoveFunc(isMn), norm.NFC)
	n, _, err := t.Transform(mnBuf, b, true)
	if err != nil {
		return nil, err
	}
	mnBuf = mnBuf[:n]
	tlBuf := bytes.NewBuffer(make([]byte, 0, len(mnBuf)*125/100))
	for i, w := 0, 0; i < len(mnBuf); i += w {
		r, width := utf8.DecodeRune(mnBuf[i:])
		if s, ok := transliterations[r]; ok {
			tlBuf.WriteString(s)
		} else {
			tlBuf.WriteRune(r)
		}
		w = width
	}
	return tlBuf.Bytes(), nil
}

func main() {
	in := "test stringß"
	fmt.Println(in)
	inBytes := []byte(in)
	outBytes, err := RemoveAccents(inBytes)
	if err != nil {
		fmt.Println(err)
	}
	out := string(outBytes)
	fmt.Println(out)
}
