package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
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

func (self *MyConsole) printf(str string, args ...any) {
	self.WriteString(fmt.Sprintf(str, args) + "\n")
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
		self.printf("%s\n", self.input[len("echo "):])
	} else if command == "type" {
		arg := fields[1]
		if slices.Contains(cmds, arg) {
			self.printf("%v is a shell builtin\n", arg)
		} else if path, _ := self.lookup(arg); path != "" {
			self.printf("%s is %s\n", arg, path)
		} else {
			self.printf("%v: not found\n", arg)
		}
	} else {
		_, err := self.lookup(command)
		if err != nil {
			self.printf("%s: command not found\n", self.input)
		} else {
			cmd := exec.Command(command, fields[1:]...)
			cmd.Stdout = os.Stdout
			cmd.Run()
		}
	}
}
