package main

import (
	"bufio"
	"fmt"
	"os"
	"slices"
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

func (self *MyConsole) Start() {
	for {
		self.WritePrompt()
		self.StoreInput()
		self.HandleInput()
	}
}

var cmds = []string{"echo", "type", "exit"}

func (self *MyConsole) HandleInput() {
	fields := strings.Fields(self.input)
	if fields[0] == "exit" {
		os.Exit(0)
	} else if fields[0] == "echo" {
		self.PrintNow(self.input[len("echo "):] + "\n")
	} else if fields[0] == "type" {
		if slices.Contains(cmds, fields[1]) {
			self.PrintNow(fmt.Sprintf("%v is a shell builtin\n", fields[1]))
		} else {
			self.PrintNow(fmt.Sprintf("%v: not found\n", fields[1]))
		}
	} else {
		self.PrintNow(fmt.Sprintf("%v: command not found\n", self.input))
	}
}
