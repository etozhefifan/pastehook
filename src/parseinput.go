package main

import (
	"os"
)


func inputArgs() (string, string) {
	fileToOpen := os.Args[1]
	if len(os.Args) <= 2 {
		return fileToOpen, nilInput
	}
	linesToSend := os.Args[2]
	return fileToOpen, linesToSend
}


