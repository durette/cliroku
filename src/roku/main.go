package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
)

func DisableInputBuffering() {
	exec.Command("stty", "-F", "/dev/tty", "cbreak", "min", "1").Run()
}

func DisableTerminalEcho() {
	exec.Command("stty", "-F", "/dev/tty", "-echo").Run()
}

func EnableTerminalEcho() {
	exec.Command("stty", "-F", "/dev/tty", "echo").Run()
}

func RokuPress(host string, key string) {
	var buf io.Reader
	resp, err := http.Post("http://"+host+":8060/keypress/"+key, "image/txt", buf)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Print(".")
	defer resp.Body.Close()
}

func SetUpControlCTrap() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Println("Done.")
		EnableTerminalEcho()
		os.Exit(0)
	}()
}

func main() {
	SetUpControlCTrap()

	var host string = "192.168.1.149"
	DisableInputBuffering()
	DisableTerminalEcho()

	fmt.Println("Arrows navigate. Enter selects. Backspace goes back. H goes home. Q or ^C quits. ")
	var b []byte = make([]byte, 1)

	for b[0] != byte('q') && b[0] != byte('Q') {
		os.Stdin.Read(b)
		switch b[0] {
		case byte(27):
			os.Stdin.Read(b)
			if b[0] == byte('[') {
				os.Stdin.Read(b)
				switch b[0] {
				case byte('A'):
					RokuPress(host, "Up")
				case byte('B'):
					RokuPress(host, "Down")
				case byte('C'):
					RokuPress(host, "Right")
				case byte('D'):
					RokuPress(host, "Left")
				}
			}
		case byte(127):
			RokuPress(host, "Back")
		case byte('h'):
			RokuPress(host, "Home")
		case byte('H'):
			RokuPress(host, "Home")
		case byte(10):
			RokuPress(host, "Select")
		default:
			fmt.Println("I got the byte", b, "("+string(b)+")")
		}
	}
	EnableTerminalEcho()
}
