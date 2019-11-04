package main

import (
	"./modules"
	"bufio"
	"fmt"
	"os"
	"path/filepath"
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

func writeCodeOneVm(commands []string, vmFilePath string, fileName string, codeWriter *modules.CodeWriter) {
	parser := modules.Parser{Commands: commands, CurrentCommand: ""}
	fileNameSplitedSlash := strings.Split(fileName, "/")
	fileName = strings.TrimSuffix(fileNameSplitedSlash[len(fileNameSplitedSlash)-1], ".vm")
	codeWriter.SetFileName(strings.TrimSuffix(fileName, ".vm"))
	codeWriter.WriteInit()
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
			codeWriter.WriteLabel(parser.Arg1())
		case "C_GOTO":
			codeWriter.WriteGoto(parser.Arg1())
		case "C_IF":
			codeWriter.WriteIf(parser.Arg1())
		case "C_FUNCTION":
			codeWriter.WriteFunction(parser.Arg1(), parser.Arg2())
		case "C_CALL":
			codeWriter.WriteCall(parser.Arg1(), parser.Arg2())
		case "C_RETURN":
			codeWriter.WriteReturn()
		default:
		}
	}
}

func main() {
	codeWriter := modules.NewCodeWriter(os.Args[1])

	err := filepath.Walk(os.Args[1], func(p string, info os.FileInfo, err error) error {
		if strings.HasSuffix(p, ".vm") {
			writeCodeOneVm(generateCommandsFromFile(p), os.Args[1], p, codeWriter)
		}
		return nil
	})
	if err != nil {
		fmt.Println(err)
	}

	codeWriter.Close() // Close output file.
}
