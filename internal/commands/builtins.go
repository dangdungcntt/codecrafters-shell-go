package commands

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func InitBuiltinCommands() {
	RegisterCommand("echo", Echo)
	RegisterCommand("type", Type)
	RegisterCommand("exit", Exit)
	RegisterCommand("pwd", Pwd)
}

func Exit(args []string) {
	code, err := strconv.Atoi(args[0])
	assertNoError(err)
	os.Exit(code)
}

func Pwd(_ []string) {
	dir, _ := os.Getwd()
	writeToConsole(dir)
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
