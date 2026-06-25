package main

import (
	"bufio"
	"fmt"
	"os"
	"unicode"

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

func (self *MyConsole) Start() {
	for {
		self.display.writePrompt()
		self.handleInput()
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

func (self *MyConsole) handleInput() {
	done := false
	mgr.resetHistoryPtr() // this needs fixing
	for !done {
		r, _, err := self.display.ReadRune()
		if err != nil || r == unicode.ReplacementChar {
			panic(err)
		} else if r == 3 {
			quitConsole(self)
		} else if r == 27 {
			captureEscapeSequence(self)
			onEscapeSequence(self)
		} else if r == 4 {
			self.display.ClearBuffer()
		} else if r == '\r' || r == '\n' {
			onReturn(self)
			done = true
		} else if r == '\t' {
			autoCompleteOnTab(self.display)
		} else if r == '\b' || r == '\x7f' {
			self.display.backspace()
		} else {
			self.display.AppendBuffer(string(r))
		}
		self.display.lastKey = r
	}
}
