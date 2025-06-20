package commands

var _ CommandInterface = Type{}

type Type struct {
	args string
}

func NewType(args string) CommandInterface {
	return Type{
		args: args,
	}
}

func (e Type) Execute() {
	_, found := CommandMap[e.args]
	if found && e.args != "type" {
		writeToConsole(e.args + " is a shell builtin")
		return
	}

	writeToConsole(e.args + ": not found")
}
