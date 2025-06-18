package commands

import (
	"fmt"
	"os"
	"strings"
)

type Command struct {
	raw string
}

func NewCommand(raw string) Command {
	return Command{
		raw: strings.TrimSuffix(raw, "\n"),
	}
}

func (c Command) Execute() {
	switch c.raw {
	case "exit 0":
		os.Exit(0)
	default:
		fmt.Println(c.raw + ": command not found")
	}
}
