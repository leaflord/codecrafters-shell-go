package main

import (
	"strconv"
)

type HistoryManager struct {
	commandHistory []string
	historyPtr     int
	escapeSequence string
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
	if mgr.historyPtr > 0 {
		mgr.historyPtr--
	}
	return mgr.commandHistory[mgr.historyPtr]
}

func (mgr *HistoryManager) getNextHistoryItem() string {
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

var mgr = &HistoryManager{[]string{}, 0, ""}

func onEscapeSequence(escapeSequence string, self *MyConsole) {
	// fmt.Println(escapeSequence)
	line := self.display.buffer
	if escapeSequence == "[A" {
		line = mgr.getLastHistoryItem()
	} else if escapeSequence == "[B" {
		line = mgr.getNextHistoryItem()
	}
	self.display.SetBuffer(line)
	// self.display.Reprompt()
	escapeSequence = ""
}
