package main

import (
	"bufio"
	"fmt"
	"os"

	"golang.org/x/term"
)

type MyConsole struct {
	fd       int
	oldState *term.State
	display  *DisplayWriter
}

const DefaultPrompt = "$ "

func NewConsole() (result *MyConsole) {
	var console = bufio.NewReadWriter(
		bufio.NewReader(os.Stdin),
		bufio.NewWriter(os.Stdout),
	)
	var display = &DisplayWriter{console, 0, "", DefaultPrompt}
	result = &MyConsole{0, nil, display}
	result.Init()
	return
}

func (self *MyConsole) Start(handler func(rune) bool) {
	for {
		self.display.writePrompt()
		self.display.handleInput(handler)
	}
}

func (self *MyConsole) Init() {
	self.fd = int(os.Stdin.Fd())
	if !term.IsTerminal(self.fd) {
		fmt.Println("Error: This must be run in a fully interactive terminal.")
		return
	}
	oldState, err := term.MakeRaw(self.fd)
	if err != nil {
		fmt.Printf("Error entering raw mode: %v\n", err)
		return
	}
	self.oldState = oldState
}

func (self *MyConsole) Clean() {
	term.Restore(self.fd, self.oldState)
}
