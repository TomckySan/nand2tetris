package main

import (
	"./assembler"
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func generateCommandsFromFile(filePath string) []string {
	commands := []string{}

	file, err := os.Open(filePath)
	if err != nil {
		// TODO
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		commands = append(commands, strings.TrimSpace(scanner.Text()))
	}

	return commands
}

func main() {
	commands := generateCommandsFromFile(os.Args[1])
	parser := assembler.Parser{Commands: commands, CurrentCommand: ""}
	code := new(assembler.Code)
	commandType := ""

	asmFileNameTrimExtension := strings.Trim(os.Args[1], ".asm")
	hackFileName := fmt.Sprintf("%s.hack", asmFileNameTrimExtension)
	hackFile, hackFileOpenErr := os.Create(hackFileName)
	if hackFileOpenErr != nil {
		// TODO
	}
	defer hackFile.Close()

	hackFileLine := ""

	for parser.HasMoreCommands() {
		parser.Advance()
		commandType = parser.CommandType()

		switch commandType {
		case "A_COMMAND":
			hackFileLine = "0"

			symbol := parser.Symbol()
			if symbolInt, err := strconv.Atoi(symbol); err == nil {
				hackFileLine += fmt.Sprintf("%015b", symbolInt)
				fmt.Fprintln(hackFile, hackFileLine)
			} else {
				// TODO
			}

		case "C_COMMAND":
			hackFileLine = "111"

			binary, err := code.Comp(parser.Comp())
			if err != nil {
				fmt.Println(err)
			}
			hackFileLine += binary

			binary, err = code.Dest(parser.Dest())
			if err != nil {
				fmt.Println(err)
			}
			hackFileLine += binary

			binary, err = code.Jump(parser.Jump())
			if err != nil {
				fmt.Println(err)
			}
			hackFileLine += binary

			fmt.Fprintln(hackFile, hackFileLine)

		case "L_COMMAND":
			// TODO
		default:
			// TODO
		}
	}
}
