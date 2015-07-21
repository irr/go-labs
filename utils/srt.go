package main

import (
	"fmt"
	"github.com/sloonz/go-iconv"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

func showInfo(s string) {
	fmt.Print(s)
	cmd := exec.Command("zenity", "--info", "--text="+s)
	cmd.Run()
}

func showError(s string) {
	cmd := exec.Command("zenity", "--error", "--text="+s)
	cmd.Run()
	log.Fatal(s)
}

func main() {

	if len(os.Args) != 2 {
		fmt.Println("usage: srt <filename>")
		os.Exit(1)
	}

	cmd := exec.Command("file", "-bi", os.Args[1])

	out, err := cmd.Output()
	if err != nil {
		showError(fmt.Sprintf("[%s] error=[%s].\n", os.Args[1], err))
	}

	charset := strings.Split(string(out), "=")
	if len(charset) != 2 {
		showError(fmt.Sprintf("[%s] charset error=[%s].\n", os.Args[1], charset))
	}

	content, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		showError(fmt.Sprintf("[%s] error=[%s].\n", os.Args[1], err))
	}

	text, err := iconv.Conv(string(content), "UTF-8", charset[1])
	if err != nil {
		text = string(content)
	}

	sub := regexp.MustCompile("(?i)<.*?i>|(?i)<.*?b>|(?i)<.*?u>").ReplaceAllString(text, "")

	if err := ioutil.WriteFile(os.Args[1], []byte(sub), 0644); err != nil {
		e := fmt.Sprintf("[%s] error=[%s].\n", os.Args[1], err)
		showError(e)
	}

	showInfo(fmt.Sprintf("[%s](%s) ok.\n", os.Args[1], charset))
}
