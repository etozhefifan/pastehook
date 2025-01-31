package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path"
)

const (
	configPath = "/.config/pastehook"
)

func checkConfig() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal("HomeDir not found")
	}
	end_path := homeDir + configPath
	found, err := os.Stat(end_path)
	if err != nil {
		fmt.Println("Dir .config not found in user directory. Would you like to initializa one? Y/n")
		configWizard(homeDir)
	}
	fmt.Println(found)
}

func configWizard(homeDir string) {
	reader := bufio.NewReader(os.Stdin)
	userInput, _, err := reader.ReadRune()
	if err != nil {
		panic("Could not get the answer")
	}
	switch userInput {
	case 'Y', 'y':
		fmt.Println("The location of the config would be", path.Join(homeDir, configPath))
		break
	case 'N', 'n':
		fmt.Println("nasss")
		break
	}
}
