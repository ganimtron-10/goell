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
	CD   = "CD"
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
	commandList := []string{EXIT, ECHO, TYPE, PWD, CD}
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

func cd(args []string) {
	if len(args) < 1 {
		return
	}

	if args[0] == "~" {
		args[0] = os.Getenv("HOME")
	}

	err := os.Chdir(args[0])
	if err != nil {
		fmt.Printf("cd: %s: No such file or directory\n", args[0])
	}
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

func modifiedFields(str string) []string {
	modFields := []string{}
	orgFields := strings.Fields(str)
	for _, ele := range orgFields {
		modFields = append(modFields, strings.ReplaceAll(ele, "\\", ""))
	}
	return modFields
}

func parseCommand(command string) []string {
	parsedCommand := []string{}

	firstSpaceLocation := strings.Index(command, " ")
	if firstSpaceLocation != -1 {
		parsedCommand = append(parsedCommand, command[:firstSpaceLocation])
	}
	command = command[firstSpaceLocation+1:]

	for {
		start := strings.IndexAny(command, "'\"")
		if start == -1 {
			parsedCommand = append(parsedCommand, modifiedFields(command)...)
			break
		}
		parsedCommand = append(parsedCommand, modifiedFields(command[:start])...)

		quote := command[start]
		command = command[start+1:]
		end := strings.IndexByte(command, quote)
		parsedCommand = append(parsedCommand, command[:end])
		command = command[end+1:]
	}

	return parsedCommand
}

func parseCommand2(command string) []string {
	parsedCommand := []string{}
	isSingleQuote := false
	isDoubleQuote := false
	// isEscaped := false
	curToken := ""

	for i := 0; i < len(command); i++ {
		char := string(command[i])

		if char == "'" {
			if isDoubleQuote {
				curToken += char
			}
			isSingleQuote = !isSingleQuote
		} else if char == "\"" {
			if isDoubleQuote {
				isDoubleQuote = false
			} else {
				if isSingleQuote {
					curToken += char
				}
				isDoubleQuote = true
			}
		} else if char == "\\" {
			if isSingleQuote {
				curToken += string(command[i])
			} else {
				if i+1 < len(command) {
					i++
					curToken += string(command[i])
				}
			}
		} else if char == " " && !isSingleQuote && !isDoubleQuote {
			if curToken != "" {
				parsedCommand = append(parsedCommand, curToken)
				curToken = ""
			}
		} else {
			curToken += char
		}

	}
	if curToken != "" {
		parsedCommand = append(parsedCommand, curToken)
	}

	// fmt.Println(parsedCommand)
	return parsedCommand
}

func parseCommand3(command string) []string {
	parsedCommand := []string{}
	isSingleQuote := false
	isDoubleQuote := false
	// isEscaped := false
	curToken := ""

	for i := 0; i < len(command); i++ {
		char := string(command[i])

		if char == "\\" && !isSingleQuote && !isDoubleQuote {
			if i+1 < len(command) {
				i++
				curToken += string(command[i])
			}
		} else if char == "\\" && isDoubleQuote {
			if i+1 < len(command) && (command[i+1] == '$' || command[i+1] == '\\' || command[i+1] == '"') {
				curToken += string(command[i+1])
				i++
			} else {
				curToken += "\\"
			}
		} else if char == "'" && !isDoubleQuote {
			isSingleQuote = !isSingleQuote
		} else if char == "\"" && !isSingleQuote {
			isDoubleQuote = !isDoubleQuote
		} else if char == " " && !isSingleQuote && !isDoubleQuote {
			if curToken != "" {
				parsedCommand = append(parsedCommand, curToken)
				curToken = ""
			}
		} else {
			curToken += char
		}

	}
	if curToken != "" {
		parsedCommand = append(parsedCommand, curToken)
		curToken = ""
	}

	// fmt.Println(parsedCommand)
	return parsedCommand
}

func evalCommand(command string) {
	// trimming new line at the end
	command = command[:len(command)-1]

	if len(command) == 0 {
		return
	}

	splittedCommand := parseCommand3(command)

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
	case CD:
		cd(splittedCommand[1:])
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
