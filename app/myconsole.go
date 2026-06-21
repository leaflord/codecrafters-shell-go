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
	buffer   string
	fd       int
	oldState *term.State
	lastKey  rune
	prompt   string
}

const DefaultPrompt = "$ "

func NewConsole() *MyConsole {
	var console = bufio.NewReadWriter(
		bufio.NewReader(os.Stdin),
		bufio.NewWriter(os.Stdout),
	)
	result := &MyConsole{console, "", 0, nil, 0, DefaultPrompt}
	result.Init()
	return result
}

func (self *MyConsole) Start() {
	for {
		self.writePrompt()
		self.onKeybuffer()
	}
}

func (self *MyConsole) writePrompt() {
	self.printNow(self.prompt)
}

func (self *MyConsole) onKeybuffer() {

	for {
		r, _, err := self.Reader.ReadRune()
		if err != nil || r == unicode.ReplacementChar {
			panic(err)
		}
		if r == 3 {
			self.Quit()
		} else if r == '\r' || r == '\n' {
			self.onReturn()
			break
		} else if r == '\t' {
			self.autoComplete()
		} else if r == '\b' || r == '\x7f' {
			self.printNow("\b \b")
			self.buffer = self.buffer[0 : len(self.buffer)-1]
		} else {
			self.AppendBuffer(string(r))
		}
	}
}

func (self *MyConsole) AppendBuffer(in string) {
	self.printNow(in)
	self.buffer = self.buffer + in
}

func (self *MyConsole) autoComplete() {
	matches := self.find(self.buffer)
	prefix := ""
	ding := len(matches) == 0

	if len(matches) > 0 {
		prefix = findCommonPrefix(matches)
		if self.buffer != prefix {
			self.AppendBuffer(prefix[len(self.buffer):] + " ")
		} else if self.lastKey == 0 {
			ding = true
			self.lastKey = '\t'
		} else {
			slices.Sort(matches)
			self.println("\r\n" + strings.Join(matches, "  "))
			self.prompt = DefaultPrompt + self.buffer
			self.writePrompt()
		}
	}
	if ding {
		self.printNow("\x07")
	}
}

func (self *MyConsole) onReturn() {
	self.buffer = strings.TrimSpace(self.buffer)
	self.printNow("\r\n")
	fields := strings.Fields(self.buffer)
	if len(fields) == 0 {
		return
	}
	command := fields[0]
	if command == "exit" {
		self.Quit()
	} else if command == "echo" {
		if len(fields) > 1 {
			self.println("%s", self.buffer[len("echo "):])
		} else {
			self.println("")
		}
	} else if command == "type" {
		arg := fields[1]
		if slices.Contains(builtinCommands, arg) {
			self.println("%v is a shell builtin", arg)
		} else if path, _ := self.lookup(arg); path != "" {
			self.println("%s is %s", arg, path)
		} else {
			self.println("%v: not found", arg)
		}
	} else {
		_, err := self.lookup(command)
		if err != nil {
			self.println("%s: command not found", self.buffer)
		} else {
			self.Clean()
			cmd := exec.Command(command, fields[1:]...)
			cmd.Stdout = os.Stdout
			cmd.Run()
			self.Init()
		}
	}
	self.buffer = ""
	self.prompt = DefaultPrompt
	self.lastKey = 0
}
