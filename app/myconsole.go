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

func (self *MyConsole) Start() {
	for {
		self.writePrompt()
		self.storeInput()
		self.handleInput()
	}
}

func (self *MyConsole) printNow(str string) {
	self.WriteString(str)
	self.Flush()
}

func (self *MyConsole) writePrompt() {
	self.printNow("$ ")
}

func (self *MyConsole) storeInput() {
	inBytes, err := self.ReadString('\n')
	if err != nil {
		panic(err)
	}
	self.input = strings.TrimSpace(inBytes)
}

var cmds = []string{"echo", "type", "exit"}

func (self *MyConsole) handleInput() {
	fields := strings.Fields(self.input)
	command := fields[0]
	if command == "exit" {
		os.Exit(0)
	} else if command == "echo" {
		self.printNow(self.input[len("echo "):] + "\n")
	} else if command == "type" {
		arg := fields[1]
		if slices.Contains(cmds, arg) {
			self.printNow(fmt.Sprintf("%v is a shell builtin\n", arg))
		} else if path, _ := self.lookup(arg); path != "" {
			self.printNow(fmt.Sprintf("%s is %s\n", arg, path))
		} else {
			self.printNow(fmt.Sprintf("%v: not found\n", arg))
		}
	} else {
		self.printNow(fmt.Sprintf("%v: command not found\n", self.input))
	}
}
