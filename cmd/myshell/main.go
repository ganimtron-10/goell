package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func evalExit(input []string) {
	exitCode, err := strconv.Atoi(input[1])

	if err != nil {
		fmt.Println("Unable to parse Exit Code. Error Details: " + err.Error())
	}

	os.Exit(exitCode)
}

func evalCommand(command string) {
	// trimming new line at the end
	command = command[:len(command)-1]

	splittedCommand := strings.Split(command, " ")
	switch splittedCommand[0] {
	case "exit":
		evalExit(splittedCommand)
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
