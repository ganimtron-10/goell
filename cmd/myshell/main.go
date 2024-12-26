package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func exit(input []string) {
	exitCode, err := strconv.Atoi(input[1])

	if err != nil {
		fmt.Println("Unable to parse Exit Code. Error Details: " + err.Error())
		os.Exit(1)
	}

	os.Exit(exitCode)
}

func echo(input []string) {
	for _, ele := range input[1:] {
		fmt.Print(ele, " ")
	}
	fmt.Println()
}

func evalCommand(command string) {
	// trimming new line at the end
	command = command[:len(command)-1]

	splittedCommand := strings.Split(command, " ")
	switch splittedCommand[0] {
	case "exit":
		exit(splittedCommand)
	case "echo":
		echo(splittedCommand)
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
