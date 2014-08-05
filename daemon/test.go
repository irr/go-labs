package main

import (
	"fmt"
	"os"
	"syscall"
	"time"
)

func daemon(nochdir, noclose int) int {
	ret, _, err := syscall.Syscall(syscall.SYS_FORK, 0, 0, 0)

	if err != 0 {
		return -1
	}

	if ret != 0 {
		os.Exit(0)
	}

	if pid, err := syscall.Setsid(); pid == -1 || err != nil {
		return -1
	}

	if nochdir == 0 {
		os.Chdir("/")
	}

	if noclose == 0 {
		f, e := os.OpenFile("/dev/null", os.O_RDWR, 0)
		if e == nil {
			fd := int(f.Fd())
			syscall.Dup2(fd, int(os.Stdin.Fd()))
			syscall.Dup2(fd, int(os.Stdout.Fd()))
			syscall.Dup2(fd, int(os.Stderr.Fd()))
		}
	}

	return 0
}

func main() {
	daemon(1, 1)
	for i := 0; i < 100000; i++ {
		fmt.Println(i)
		time.Sleep(1 * time.Second)
	}
}
