package main

import (
	"maps"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"
)

var executableFiles = []string{}
var commands = []string{"echo", "type", "exit", "history"}

var commandHistory = []string{}

func init() {
	result := make(map[string](struct{}))

	for _, builtin := range commands {
		result[builtin] = struct{}{}
	}

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
	slices.Sort(executableFiles)
}

func lookup(fileName string) (out string, err error) {
	return exec.LookPath(fileName)
}

func quitConsole(self *MyConsole) {
	self.Clean()
	os.Exit(0)
}

func find(filePrefix string) []string {
	result := make([]string, 0)
	for _, f := range executableFiles {
		if strings.HasPrefix(f, filePrefix) {
			result = append(result, f)
		}
	}
	return result
}

func findCommonPrefix(strings []string) (result string) {
	if len(strings) == 0 {
		return
	}
	runeseqs := [][]rune{}
	for _, str := range strings {
		runeseqs = append(runeseqs, []rune(str))
	}

	for matches, i, fst := true, 0, runeseqs[0]; i < len(fst) && matches; i++ {
		for j := 1; j < len(runeseqs) && matches; j++ {
			matches = len(runeseqs[j]) < i || fst[i] == runeseqs[j][i]
		}
		if matches {
			result = result + string(fst[i])
		}
	}
	return
}

func handleRune(self *MyConsole, r rune) bool {
	if r == 3 {
		quitConsole(self)
	} else if r == 4 {
		self.display.ClearBuffer()
	} else if r == '\r' || r == '\n' {
		onReturn(self)
		return true
	} else if r == '\t' {
		autoCompleteOnTab(self.display)
	} else if r == '\b' || r == '\x7f' {
		self.display.backspace()
	} else {
		self.display.AppendBuffer(string(r))
	}
	return false
}

func autoCompleteOnTab(self *DisplayWriter) {
	// to clean further, "findCommonPrefix" can be compared with existing buffer beforehand
	matches := find(self.buffer)
	if len(matches) == 0 { // no matches
		self.printNow("\x07")
		return
	}

	prefix := findCommonPrefix(matches)
	self.SetBuffer(prefix)
	if len(matches) == 1 { // single match
		self.AppendBuffer(" ")
	} else if prefix != self.buffer { // multi-matches with shared prefix
		self.Reprompt()
	} else if self.lastKey != '\t' { // multi-match without shared prefix, tab pressed once
		self.printNow("\x07")
	} else { // multi-match without shared prefix, tab pressed twice
		slices.Sort(matches)
		self.println("\r\n" + strings.Join(matches, "  "))
		self.Reprompt()
	}
}

func onReturn(self *MyConsole) {
	display := self.display
	inputline := strings.TrimSpace(display.buffer)
	commandHistory = append(commandHistory, inputline)
	display.printNow("\r\n")
	fields := strings.Fields(inputline)
	if len(fields) == 0 {
		return
	}
	command := fields[0]
	if command == "echo" {
		display.println(strings.Join(fields[1:], " "))
	} else if command == "type" {
		arg := fields[1]
		if slices.Contains(commands, arg) {
			display.println("%v is a shell builtin", arg)
		} else if path, _ := lookup(arg); path != "" {
			display.println("%s is %s", arg, path)
		} else {
			display.println("%v: not found", arg)
		}
	} else if command == "history" {
		showHistory(self)
	} else if command == "exit" {
		quitConsole(self)
	} else {
		executeCommand(self, fields)
	}
	display.Reset()
}

func executeCommand(console *MyConsole, fields []string) {
	_, err := lookup(fields[0])
	if err != nil {
		console.display.println("%s: command not found", fields[0])
		return
	}
	console.Clean()
	defer console.Init()
	cmd := exec.Command(fields[0], fields[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func showHistory(self *MyConsole) {
	for i, line := range commandHistory {
		self.display.println("\t%v %s", 1+i, line)
	}
}
