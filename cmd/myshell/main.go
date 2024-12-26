package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

const (
	EXIT = "EXIT"
	ECHO = "ECHO"
	TYPE = "TYPE"
)

func exit(args []string) {
	exitCode, err := strconv.Atoi(args[0])

	if err != nil {
		fmt.Println("Unable to parse Exit Code. Error Details: " + err.Error())
		os.Exit(1)
	}

	os.Exit(exitCode)
}

func echo(args []string) {
	for _, ele := range args {
		fmt.Print(ele, " ")
	}
	fmt.Println()
}

func evalType(args []string) {
	commandList := []string{EXIT, ECHO, TYPE}
	isInbuilt := false
	executablePath := ""

	// check if builtin
	for _, ele := range commandList {
		if strings.ToUpper(args[0]) == ele {
			isInbuilt = true
			break
		}
	}

	// check if present in PATH
	pathValue := os.Getenv("PATH")
	pathList := strings.Split(pathValue, ":")
	for _, path := range pathList {
		fullPath := filepath.Join(path, args[0])

		_, err := os.Stat(fullPath)
		if err == nil {
			executablePath = fullPath
			break
		}

		if !os.IsNotExist(err) {
			fmt.Printf("Error while checking path %s. Error Details: %s\n", fullPath, err.Error())
			os.Exit(1)
		}

	}

	if isInbuilt {
		fmt.Println(args[0] + " is a shell builtin")
	} else if executablePath != "" {
		fmt.Println(args[0] + " is " + executablePath)
	} else {
		fmt.Println(args[0] + ": not found")
	}
}

func evalCommand(command string) {
	// trimming new line at the end
	command = command[:len(command)-1]

	if len(command) == 0 {
		return
	}

	splittedCommand := strings.Split(command, " ")

	// if len(splittedCommand) < 2 {
	// 	fmt.Printf("No args provided for command %s\n", splittedCommand[0])
	// 	return
	// }

	switch strings.ToUpper(splittedCommand[0]) {
	case EXIT:
		exit(splittedCommand[1:])
	case ECHO:
		echo(splittedCommand[1:])
	case TYPE:
		evalType(splittedCommand[1:])
	default:
		fmt.Println(command + ": command not found")
	}
}

func main() {

	for {
		fmt.Fprint(os.Stdout, "$ ")

		// Wait for user input
		command, err := bufio.NewReader(os.Stdin).ReadString('\n')

		if err != nil {
			fmt.Fprintln(os.Stderr, "Error reading input:", err)
			os.Exit(1)
		}

		evalCommand(command)
	}

}
