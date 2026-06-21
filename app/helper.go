package main

import (
	"fmt"
	"maps"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"

	"golang.org/x/term"
)

var executableFiles = []string{}
var builtinCommands = []string{"echo", "type", "exit"}

func init() {
	result := make(map[string](struct{}))
	pathlist := os.Getenv("PATH")
	paths := strings.Split(pathlist, ":")
	for _, path := range paths {
		files, err := os.ReadDir(path)
		if err != nil || files == nil {
			continue // graceful return
		}
		for _, file := range files {
			absPath, _ := filepath.Abs(filepath.Join(path, file.Name()))
			if err != nil {
				continue
			}
			stat, err := os.Stat(absPath)
			if err != nil {
				continue
			}
			if (stat.Mode().Perm() & 0100) != 0 {
				if _, ok := result[file.Name()]; !ok {
					result[file.Name()] = struct{}{}
				}
			}
		}
	}
	executableFiles = slices.Collect(maps.Keys(result))
	// executableFiles = append(executableFiles, builtinCommands...)
	slices.Sort(executableFiles)
}

func (self *MyConsole) lookup(fileName string) (out string, err error) {
	return exec.LookPath(fileName)
}

func (self *MyConsole) printNow(str string) {
	self.WriteString(str)
	self.Flush()
}

func (self *MyConsole) println(str string, args ...any) {
	self.WriteString(fmt.Sprintf(str, args...))
	self.WriteString("\r\n")
	self.Flush()
}

func (self *MyConsole) find(filePrefix string) []string {
	result := make([]string, 0)
	for _, f := range builtinCommands {
		if strings.HasPrefix(f, filePrefix) {
			result = append(result, f)
		}
	}
	return result
}

func (self *MyConsole) Quit() {
	self.Clean()
	os.Exit(0)
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
