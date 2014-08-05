package main

import (
	"fmt"
	"github.com/sloonz/go-iconv"
	"io/ioutil"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

func main() {

	if len(os.Args) != 2 {
		fmt.Println("usage: srt <filename>")
		os.Exit(1)
	}

	cmd := exec.Command("file", "-bi", os.Args[1])

	out, err := cmd.Output()
	if err != nil {
		os.Exit(1)
	}

	charset := strings.Split(string(out), "=")
	if len(charset) != 2 {
		os.Exit(1)
	}

	content, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		os.Exit(1)
	}

	text, err := iconv.Conv(string(content), "UTF-8", charset[1])
	if err != nil {
		text = string(content)
	}

	sub := regexp.MustCompile("(?i)<.*?i>|(?i)<.*?b>|(?i)<.*?u>").ReplaceAllString(text, "")

	if err := ioutil.WriteFile(os.Args[1], []byte(sub), 0644); err != nil {
		os.Exit(1)
	}
}
