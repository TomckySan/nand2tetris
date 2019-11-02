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
	// TODO
	switch {
	case p.isArithmetic(p.CurrentCommand):
		return "C_ARITHMETIC"
	default:
		return "C_PUSH"
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
		// TODO
		return ""
	case "C_LABEL":
		// TODO
		return ""
	case "C_GOTO":
		// TODO
		return ""
	case "C_IF":
		// TODO
		return ""
	case "C_FUNCTION":
		// TODO
		return ""
	case "C_RETURN":
		// TODO
		return ""
	case "C_CALL":
		// TODO
		return ""
	default:
		// TODO
		return ""
	}
}

func (p Parser) Arg2() int {
	command := p.CurrentCommand
	v, _ := strconv.Atoi(strings.Fields(command)[2])
	return v
}
