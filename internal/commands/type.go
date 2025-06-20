package commands

import (
	"fmt"
	"os"
	"strings"
)

var _ CommandInterface = Type{}

type Type struct {
	args []string
}

func NewType(args []string) CommandInterface {
	return Type{
		args: args,
	}
}

func (e Type) Execute() {
	bin := e.args[0]
	_, found := CommandMap[bin]
	if found {
		writeToConsole(bin + " is a shell builtin")
		return
	}

	if file, exists := findBinInPath(bin); exists {
		writeToConsole(fmt.Sprintf("%s is %s", bin, file))
		return
	}

	writeToConsole(bin + ": not found")
}

func findBinInPath(bin string) (string, bool) {
	paths := os.Getenv("PATH")
	for _, path := range strings.Split(paths, ":") {
		file := path + "/" + bin
		if _, err := os.Stat(file); err == nil {
			return file, true
		}
	}

	return "", false
}
