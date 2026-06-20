package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type MyConsole struct {
	*bufio.ReadWriter
	input string
}

func NewConsole() *MyConsole {
	var console = bufio.NewReadWriter(
		bufio.NewReader(os.Stdin),
		bufio.NewWriter(os.Stdout),
	)
	return &MyConsole{console, ""}
}

func (self *MyConsole) PrintNow(str string) {
	self.WriteString(str)
	self.Flush()
}

func (self *MyConsole) WritePrompt() {
	self.PrintNow("$ ")
}

func (self *MyConsole) StoreInput() {
	inBytes, err := self.ReadString('\n')
	if err != nil {
		panic(err)
	}
	self.input = strings.TrimSpace(inBytes)
}

func (self *MyConsole) HandleInput() {
	if self.input == "exit" {
		os.Exit(0)
	} else {
		self.PrintNow(fmt.Sprintf("%v: command not found\n", self.input))
	}
}

func (self *MyConsole) Start() {
	for {
		self.WritePrompt()
		self.StoreInput()
		self.HandleInput()
	}
}
