package main

import (
	"bufio"
	"os"
	"os/exec"
	"slices"
	"strings"
	"unicode"

	"golang.org/x/term"
)

type MyConsole struct {
	*bufio.ReadWriter
	input    string
	fd       int
	oldState *term.State
}

func NewConsole() *MyConsole {
	var console = bufio.NewReadWriter(
		bufio.NewReader(os.Stdin),
		bufio.NewWriter(os.Stdout),
	)
	result := &MyConsole{console, "", 0, nil}
	result.Init()
	return result
}

func (self *MyConsole) Start() {
	for {
		self.writePrompt()
		self.onKeyInput()
	}
}

func (self *MyConsole) writePrompt() {
	self.printNow("$ ")
}

func (self *MyConsole) onKeyInput() {

	buffer := ""
	for {
		r, _, err := self.Reader.ReadRune()
		if err != nil || r == unicode.ReplacementChar {
			panic(err)
		}
		if r == 3 {
			self.Quit()
		} else if r == '\r' || r == '\n' {
			break
		} else if r == '\t' {
			self.autoComplete(buffer)
		} else if r == '\b' || r == '\x7f' {
			self.printNow("\b \b")
			buffer = buffer[0 : len(buffer)-1]
		} else {
			bufch := string(r)
			buffer = buffer + bufch
			self.printNow(bufch)
		}
	}

	// inBytes, err := self.ReadString('\n')
	// if err != nil {
	// 	panic(err)
	// }
	self.input = strings.TrimSpace(buffer)
	self.onReturn()
}

func (self *MyConsole) autoComplete(buffer string) {
	matches := self.find(buffer)
	if len(matches) == 1 && len(matches[0]) > len(buffer) {
		self.printNow(matches[0][len(buffer):])
	}
}

var cmds = []string{"echo", "type", "exit"}

func (self *MyConsole) onReturn() {
	self.printNow("\r\n")
	fields := strings.Fields(self.input)
	if len(fields) == 0 {
		return
	}
	command := fields[0]
	if command == "exit" {
		os.Exit(0)
	} else if command == "echo" {
		if len(fields) > 1 {
			self.println("%s", self.input[len("echo "):])
		} else {
			self.println("")
		}
	} else if command == "type" {
		arg := fields[1]
		if slices.Contains(cmds, arg) {
			self.println("%v is a shell builtin", arg)
		} else if path, _ := self.lookup(arg); path != "" {
			self.println("%s is %s", arg, path)
		} else {
			self.println("%v: not found", arg)
		}
	} else {
		_, err := self.lookup(command)
		if err != nil {
			self.println("%s: command not found", self.input)
		} else {
			cmd := exec.Command(command, fields[1:]...)
			cmd.Stdout = os.Stdout
			cmd.Run()
		}
	}
	self.input = ""
}
