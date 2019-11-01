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

// First Path: Only creating symbol table
func firstPath(commands []string, symbolTable *assembler.SymbolTable) {
	parser := assembler.Parser{Commands: commands, CurrentCommand: ""}
	romAddress := 0
	commandType := ""

	for parser.HasMoreCommands() {
		parser.Advance() // read next command
		commandType = parser.CommandType()

		if commandType == "A_COMMAND" || commandType == "C_COMMAND" {
			romAddress++
		} else if commandType == "L_COMMAND" {
			symbolTable.AddEntry(parser.Symbol(), romAddress)
		} else {
			// Ignore
		}
	}
}

// Second Path:
func secondPath(commands []string, symbolTable *assembler.SymbolTable) {
	// Open output file
	asmFileNameTrimExtension := strings.Trim(os.Args[1], ".asm")
	hackFileName := fmt.Sprintf("%s.hack", asmFileNameTrimExtension)
	hackFile, hackFileOpenErr := os.Create(hackFileName)
	if hackFileOpenErr != nil {
		fmt.Println(hackFileOpenErr)
	}
	defer hackFile.Close()

	// Definition
	parser := assembler.Parser{Commands: commands, CurrentCommand: ""}
	commandType := ""
	hackFileLine := ""
	code := new(assembler.Code)

	for parser.HasMoreCommands() {
		parser.Advance() // read next command
		commandType = parser.CommandType()

		switch commandType {
		case "A_COMMAND":
			hackFileLine = "0"
			symbol := parser.Symbol()
			if symbolInt, err := strconv.Atoi(symbol); err == nil {
				hackFileLine += fmt.Sprintf("%015b", symbolInt)
				fmt.Fprintln(hackFile, hackFileLine)
			} else if symbolTable.Contains(symbol) {
				hackFileLine += fmt.Sprintf("%015b", symbolTable.GetAddress(symbol))
				fmt.Fprintln(hackFile, hackFileLine)
			} else {
				hackFileLine += fmt.Sprintf("%015b", symbolTable.GetNextAddress())
				fmt.Fprintln(hackFile, hackFileLine)
				symbolTable.AddEntry(parser.Symbol(), symbolTable.GetNextAddress())
				symbolTable.IncrementNextAddress()
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
			// Ignore
		default:
			// Ignore
		}
	}
}

func main() {
	commands := generateCommandsFromFile(os.Args[1])
	symbolTable := assembler.NewSymbolTable()

	firstPath(commands, symbolTable)
	secondPath(commands, symbolTable)
}
