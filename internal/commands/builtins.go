package commands

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"strconv"
	"strings"
)

func Exit(args []string) {
	code, err := strconv.Atoi(args[0])
	assertNoError(err)
	os.Exit(code)
}

func Pwd(_ []string) {
	writeToConsole(State.Pwd)
}

func Cd(args []string) {
	newPath := path.Join(State.Pwd, args[0])
	if !IsExist(newPath) {
		writeToConsole(fmt.Sprintf("cd: %s: No such file or directory", args[0]))
		return
	}

	State.Pwd = newPath
}

func Echo(args []string) {
	writeToConsole(strings.Join(args, " "))
}

func Type(args []string) {
	bin := args[0]
	if IsBuiltin(bin) {
		writeToConsole(bin + " is a shell builtin")
		return
	}

	if file, exists := findBinInPath(bin); exists {
		writeToConsole(fmt.Sprintf("%s is %s", bin, file))
		return
	}

	writeToConsole(bin + ": not found")
}

func RunExternalApp(executable string, args []string) {
	_, found := findBinInPath(executable)
	if found {
		cmd := exec.Command(executable, args...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Run()
	} else {
		writeToConsole(executable + ": command not found")
	}
}
