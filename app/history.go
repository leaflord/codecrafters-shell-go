package main

import (
	"strconv"
	"strings"
)

type HistoryManager struct {
	commandHistory []string
	historyPtr     int
	escapeSequence strings.Builder
}

func (mgr *HistoryManager) showHistory(self *MyConsole, fields []string) {
	count := len(mgr.commandHistory)
	if len(fields) > 1 {
		tmp, err := strconv.Atoi(fields[1])
		if err != nil {
			panic(err)
		}
		count = tmp

	}

	for i, line := range mgr.commandHistory {
		if i >= len(mgr.commandHistory)-count {
			self.display.println("\t%v %s", 1+i, line)
		}
	}
}

func (mgr *HistoryManager) getLastHistoryItem() string {
	if len(mgr.commandHistory) == 0 {
		return ""
	}
	if mgr.historyPtr > 0 {
		mgr.historyPtr--
	}
	return mgr.commandHistory[mgr.historyPtr]
}

func (mgr *HistoryManager) getNextHistoryItem() string {
	if len(mgr.commandHistory) == 0 {
		return ""
	}
	if mgr.historyPtr < len(mgr.commandHistory)-1 {
		mgr.historyPtr++
	}
	return mgr.commandHistory[mgr.historyPtr]
}

func (mgr *HistoryManager) resetHistoryPtr() {
	mgr.historyPtr = len(mgr.commandHistory)
}

func (mgr *HistoryManager) addEntry(str string) {
	mgr.commandHistory = append(mgr.commandHistory, str)
}

var mgr = &HistoryManager{[]string{}, 0, strings.Builder{}}

func onEscapeSequence(self *MyConsole) {
	val := mgr.escapeSequence.String()
	line := self.display.buffer
	if val == "[A" {
		line = mgr.getLastHistoryItem()
	} else if val == "[B" {
		line = mgr.getNextHistoryItem()
	}
	mgr.escapeSequence.Reset()
	self.display.SetBuffer(line)
}

func captureEscapeSequence(console *MyConsole) {
	fstChar := true
	for {
		r, _, _ := console.display.ReadRune()
		mgr.escapeSequence.WriteRune(r)
		if fstChar && r == '[' {
			fstChar = false
		} else if r >= 0x40 && r <= 0x7E {
			break
		}
	}
}
