package modules

import (
	"strconv"
	"strings"
)

type Parser struct {
	Commands       []string
	CurrentCommand string
}

func (p Parser) HasMoreCommands() bool {
	return len(p.Commands) != 0
}

func (p *Parser) Advance() {
	p.CurrentCommand = p.Commands[0]
	p.Commands = p.Commands[1:]
}

func (p Parser) CommandType() string {
	command := strings.Fields(p.CurrentCommand)[0]
	switch {
	case p.isArithmetic(p.CurrentCommand):
		return "C_ARITHMETIC"
	case command == "push":
		return "C_PUSH"
	case command == "pop":
		return "C_POP"
	case command == "label":
		return "C_LABEL"
	case command == "goto":
		return "C_GOTO"
	case command == "if-goto":
		return "C_IF"
	case command == "function":
		return "C_FUNCTION"
	case command == "call":
		return "C_CALL"
	case command == "return":
		return "C_RETURN"
	default:
		return ""
	}
}

func (p Parser) isArithmetic(command string) bool {
	return strings.HasPrefix(command, "add") ||
		strings.HasPrefix(command, "sub") ||
		strings.HasPrefix(command, "neg") ||
		strings.HasPrefix(command, "eq") ||
		strings.HasPrefix(command, "gt") ||
		strings.HasPrefix(command, "lt") ||
		strings.HasPrefix(command, "and") ||
		strings.HasPrefix(command, "or") ||
		strings.HasPrefix(command, "not")
}

func (p Parser) Arg1() string {
	command := p.CurrentCommand
	commandType := p.CommandType()
	switch commandType {
	case "C_ARITHMETIC":
		return strings.Fields(command)[0]
	case "C_PUSH":
		return strings.Fields(command)[1]
	case "C_POP":
		return strings.Fields(command)[1]
	case "C_LABEL":
		return strings.Fields(command)[1]
	case "C_GOTO":
		return strings.Fields(command)[1]
	case "C_IF":
		return strings.Fields(command)[1]
	case "C_FUNCTION":
		return strings.Fields(command)[1]
	case "C_CALL":
		return strings.Fields(command)[1]
	default:
		return ""
	}
}

func (p Parser) Arg2() int {
	command := p.CurrentCommand
	v, _ := strconv.Atoi(strings.Fields(command)[2])
	return v
}
