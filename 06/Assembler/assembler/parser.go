package assembler

import (
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
	switch {
	case strings.HasPrefix(p.CurrentCommand, "@"):
		return "A_COMMAND"
	case strings.HasPrefix(p.CurrentCommand, "("):
		return "L_COMMAND"
	default:
		return "C_COMMAND"
	}
}

func (p Parser) Symbol() string {
	c := p.CurrentCommand
	switch {
	case strings.HasPrefix(p.CurrentCommand, "@"):
		return strings.Trim(c, "@")
	default:
		return strings.Trim(c, "()")
	}
}

func (p Parser) Dest() string {
	c := p.CurrentCommand
	if strings.Contains(c, "=") {
		return strings.Split(c, "=")[0]
	}
	return "null"
}

func (p Parser) Comp() string {
	c := p.CurrentCommand
	if strings.Contains(c, "=") {
		return strings.Split(strings.Split(c, "=")[1], ";")[0]
	}
	return strings.Split(c, ";")[0]
}

func (p Parser) Jump() string {
	c := p.CurrentCommand
	if strings.Contains(c, ";") {
		return strings.Split(c, ";")[1]
	}
	return "null"
}
