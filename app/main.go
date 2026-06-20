package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"golang.org/x/term"
)

func main() {
	fd := int(os.Stdin.Fd())
	if !term.IsTerminal(fd) {
		fmt.Println("Error: This must be run in a fully interactive terminal.")
		return
	}

	// Turn on raw mode to block terminal echo
	// oldState, err := term.MakeRaw(fd)
	// if err != nil {
	// 	fmt.Printf("Error entering raw mode: %v\n", err)
	// 	return
	// }
	// defer term.Restore(fd, oldState)

	var console = bufio.NewReadWriter(
		bufio.NewReader(os.Stdin),
		bufio.NewWriter(os.Stdout),
	)
	start(console)
}

func start(console *bufio.ReadWriter) {
	for {
		console.WriteString("$ ")
		console.Flush()
		in, err := console.ReadString('\n')
		if err != nil {
			panic(err)
		}
		console.WriteString(fmt.Sprintf("%v: command not found", strings.TrimSpace(in)))
		console.Flush()
	}
}
