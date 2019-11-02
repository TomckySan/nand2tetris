package main

import (
	"./modules"
	"bufio"
	"fmt"
	"os"
	"strings"
)

func generateCommandsFromFile(filePath string) []string {
	commands := []string{}

	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		spaceTrimedLine := strings.TrimSpace(scanner.Text())
		if len(spaceTrimedLine) == 0 || strings.HasPrefix(spaceTrimedLine, "//") {
			continue // Comment line skip.
		}
		commands = append(commands, strings.TrimSpace(strings.Split(spaceTrimedLine, "//")[0]))
	}

	return commands
}

func main() {
	vmFilePath := os.Args[1]
	commands := generateCommandsFromFile(vmFilePath)
	parser := modules.Parser{Commands: commands, CurrentCommand: ""}
	codeWriter := modules.NewCodeWriter(vmFilePath)
	vmFilePathSplitedSlash := strings.Split(vmFilePath, "/")
	codeWriter.SetFileName(vmFilePathSplitedSlash[len(vmFilePathSplitedSlash)-1])
	commandType := ""

	for parser.HasMoreCommands() {
		parser.Advance() // read next command
		commandType = parser.CommandType()

		switch commandType {
		case "C_ARITHMETIC":
			codeWriter.WriteArithmetic(parser.Arg1())
		case "C_PUSH":
			segment := parser.Arg1()
			index := parser.Arg2()
			codeWriter.WritePushPop(commandType, segment, index)
		case "C_POP":
			segment := parser.Arg1()
			index := parser.Arg2()
			codeWriter.WritePushPop(commandType, segment, index)
		case "C_LABEL":
			// ここでは実装しない
		case "C_GOTO":
			// ここでは実装しない
		case "C_IF":
			// ここでは実装しない
		case "C_FUNCTION":
			// ここでは実装しない
		case "C_RETURN":
			// ここでは実装しない
		case "C_CALL":
			// ここでは実装しない
		default:
		}
	}

	codeWriter.Close() // Close output file.
}
