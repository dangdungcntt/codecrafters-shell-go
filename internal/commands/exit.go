package commands

import (
	"os"
	"strconv"
)

var _ CommandInterface = Exit{}

type Exit struct {
	args string
}

func NewExit(args string) CommandInterface {
	return Exit{
		args: args,
	}
}

func (e Exit) Execute() {
	code, err := strconv.Atoi(e.args)
	assertNoError(err)
	os.Exit(code)
}
