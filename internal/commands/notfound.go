package commands

import "fmt"

type NotFound struct {
	executable string
	args       []string
}

func NewNotFoundHandler(executable string, args []string) CommandInterface {
	return NotFound{
		executable: executable,
		args:       args,
	}
}

func (n NotFound) Execute() {
	fmt.Println(n.executable + ": command not found")
}
