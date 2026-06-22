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

func NewConsole() (result *MyConsole) {
	var console = bufio.NewReadWriter(
		bufio.NewReader(os.Stdin),
		bufio.NewWriter(os.Stdout),
	)
	result = &MyConsole{console, "", 0, nil, 0, DefaultPrompt}
	result.Init()
	return
}

func (self *MyConsole) Start() {
	for {
		self.writePrompt()
		self.handleInput()
	}
}

func (self *MyConsole) writePrompt() {
	self.printNow(self.prompt)
}

func (self *MyConsole) handleInput() {
	loop := true
	for loop {
		r, _, err := self.Reader.ReadRune()
		if err != nil || r == unicode.ReplacementChar {
			panic(err)
		}
		if r == 3 {
			self.Quit()
		} else if r == '\r' || r == '\n' {
			self.onReturn()
			loop = false
		} else if r == '\t' {
			self.autoCompleteOnTab()
		} else if r == '\b' || r == '\x7f' {
			self.printNow("\b \b")
			self.buffer = self.buffer[0 : len(self.buffer)-1]
		} else {
			self.AppendBuffer(string(r))
		}
		self.lastKey = r
	}
}

func (self *MyConsole) autoCompleteOnTab() {
	// to clean further, "findCommonPrefix" can be compared with existing buffer beforehand
	matches := self.find(self.buffer)
	if len(matches) == 0 { // no matches
		self.printNow("\x07")
		return
	}

	prefix := findCommonPrefix(matches)
	self.AppendBuffer(prefix[len(self.buffer):])
	if len(matches) == 1 { // single match
		self.AppendBuffer(" ")
	} else if prefix != self.buffer { // multi-matches with shared prefix
		self.prompt = DefaultPrompt + self.buffer
		self.writePrompt()
	} else if self.lastKey != '\t' { // multi-match without shared prefix, tab pressed once
		self.printNow("\x07")
	} else { // multi-match without shared prefix, tab pressed twice
		slices.Sort(matches)
		self.println("\r\n" + strings.Join(matches, "  "))
		self.prompt = DefaultPrompt + self.buffer
		self.writePrompt()
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
			defer self.Init()
			cmd := exec.Command(command, fields[1:]...)
			cmd.Stdout = os.Stdout
			cmd.Run()
		}
	}
	self.buffer = ""
	self.prompt = DefaultPrompt
}
