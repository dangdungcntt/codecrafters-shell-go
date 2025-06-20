package commands

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
	writeToConsole(e.args)
}
