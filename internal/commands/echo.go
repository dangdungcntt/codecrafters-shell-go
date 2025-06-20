package commands

import (
	"fmt"
	"os"
)

var _ CommandInterface = Echo{}

type Echo struct {
	args string
}

func NewEcho(args string) CommandInterface {
	return Echo{
		args: args,
	}
}

func (e Echo) Execute() {
	fmt.Fprint(os.Stdout, e.args, "\n")
}
