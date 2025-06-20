package commands

import (
	"fmt"
	"os"
)

type CommandInterface interface {
	Execute()
}

var CommandMap = map[string]func(args string) CommandInterface{
	"echo": NewEcho,
	"exit": NewExit,
	"type": NewType,
}

func NewCommand(executable string, args string) CommandInterface {
	constructor, found := CommandMap[executable]
	if !found {
		return NewNotFoundHandler(executable, args)
	}

	return constructor(args)
}

func assertNoError(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func writeToConsole(args ...string) {
	for _, arg := range args {
		fmt.Fprint(os.Stdout, arg)
	}

	fmt.Fprint(os.Stdout, "\n")
}
