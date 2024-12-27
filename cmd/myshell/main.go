package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

const (
	EXIT = "EXIT"
	ECHO = "ECHO"
	TYPE = "TYPE"
	PWD  = "PWD"
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

func checkExecutable(commandName string) string {
	pathValue := os.Getenv("PATH")
	pathList := strings.Split(pathValue, ":")
	for _, path := range pathList {
		fullPath := filepath.Join(path, commandName)

		_, err := os.Stat(fullPath)
		if err == nil {
			return fullPath
		}

		if !os.IsNotExist(err) {
			fmt.Printf("Error while checking path %s. Error Details: %s\n", fullPath, err.Error())
		}

	}
	return ""
}

func evalType(args []string) {
	commandList := []string{EXIT, ECHO, TYPE, PWD}
	isInbuilt := false

	// check if builtin
	for _, ele := range commandList {
		if strings.ToUpper(args[0]) == ele {
			isInbuilt = true
			break
		}
	}

	// check if present in PATH
	executablePath := checkExecutable(args[0])

	if isInbuilt {
		fmt.Println(args[0] + " is a shell builtin")
	} else if executablePath != "" {
		fmt.Println(args[0] + " is " + executablePath)
	} else {
		fmt.Println(args[0] + ": not found")
	}
}

func pwd() {
	dirPath, err := os.Getwd()
	if err != nil {
		fmt.Println("Error while getting current working dirtectory. Error details: " + err.Error())
		os.Exit(1)
	}
	fmt.Println(dirPath)
}

func execute(executablePath string, args []string) {
	command := exec.Command(executablePath, args...)
	command.Stderr = os.Stderr
	command.Stdout = os.Stdout
	err := command.Run()
	if err != nil {
		fmt.Println("Error while executing " + executablePath + ". Error details: " + err.Error())
		os.Exit(1)
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
	case PWD:
		pwd()
	default:
		if executablePath := checkExecutable(splittedCommand[0]); executablePath != "" {
			execute(executablePath, splittedCommand[1:])
		} else {
			fmt.Println(command + ": command not found")
		}

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
