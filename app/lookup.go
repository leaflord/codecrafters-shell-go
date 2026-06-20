package main

import (
	"os/exec"
)

func (self *MyConsole) lookup(fileName string) (out string, err error) {
	return exec.LookPath(fileName)
}
