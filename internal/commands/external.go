package commands

import (
	"os"
	"os/exec"
)

type External struct {
	executable string
	args       []string
}

func NewExternal(executable string, args []string) CommandInterface {
	return External{
		executable: executable,
		args:       args,
	}
}

func (n External) Execute() {
	_, found := findBinInPath(n.executable)
	if found {
		cmd := exec.Command(n.executable, n.args...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Run()
	} else {
		writeToConsole(n.executable + ": command not found")
	}
}
