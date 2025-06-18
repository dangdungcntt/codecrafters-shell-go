package commands

import (
	"errors"
	"fmt"
	"os"
	"strconv"
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

func (c Command) parse() (string, []string, error) {
	parts := strings.SplitN(c.raw, " ", 2)
	if len(parts) == 0 {
		return "", nil, errors.New("empty command")
	}

	var args []string
	if len(parts) == 2 {
		args = strings.Split(parts[1], " ")
	}

	return parts[0], args, nil
}

func (c Command) Execute() {
	executable, args, err := c.parse()
	assertNoError(err)

	switch executable {
	case "exit":
		code, err := strconv.Atoi(args[0])
		assertNoError(err)
		os.Exit(code)
	case "echo":
		fmt.Fprint(os.Stdout, args[0], "\n")
	default:
		fmt.Println(c.raw + ": command not found")
	}
}

func assertNoError(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
