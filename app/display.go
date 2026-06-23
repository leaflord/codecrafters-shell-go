package main

import (
	"bufio"
	"fmt"
)

type DisplayWriter struct {
	*bufio.ReadWriter
	lastKey rune
	buffer  string
	prompt  string
}

func (self *DisplayWriter) writePrompt() {
	self.printNow(self.prompt)
}

func (self *DisplayWriter) printNow(str string) {
	self.WriteString(str)
	self.Flush()
}

func (self *DisplayWriter) println(str string, args ...any) {
	self.WriteString(fmt.Sprintf(str, args...))
	self.WriteString("\r\n")
	self.Flush()
}

func (self *DisplayWriter) AppendBuffer(in string) {
	self.printNow(in)
	self.buffer = self.buffer + in
}

func (self *DisplayWriter) backspace() {
	self.printNow("\b \b")
	self.buffer = self.buffer[0 : len(self.buffer)-1]
}

func (self *DisplayWriter) ClearBuffer() {
	buf := ""
	for range len(self.buffer) {
		buf = buf + "\b \b"
	}
	self.printNow(buf)
	self.Reset()
}

func (self *DisplayWriter) SetBuffer(in string) {
	self.ClearBuffer()
	self.AppendBuffer(in)
	self.buffer = in
}

func (self *DisplayWriter) Reprompt() {
	self.prompt = DefaultPrompt + self.buffer
	self.writePrompt()
}

func (self *DisplayWriter) Reset() {
	self.buffer = ""
	self.prompt = DefaultPrompt
}
