package main

import (
	"flag"
)

func inputArgs() (string, string) {
	var fileToOpen string
	var linesToSend string
	flag.StringVar(&fileToOpen, "f", "", "specify a file for gathering lines from it")
	flag.StringVar(&linesToSend, "l", "", "specify lines to send from file, divided by -. If omitted, the whole file will be sent to pastebin.")
	flag.Parse()
	return fileToOpen, linesToSend
}

