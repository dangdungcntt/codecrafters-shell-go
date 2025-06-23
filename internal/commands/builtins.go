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
	writeOutput(State.Cwd())
}

func Cd(args []string) {
	var targetPath string
	switch {
	case len(args) == 0 || strings.HasPrefix(args[0], "~"):
		targetPath, _ = os.UserHomeDir()
		if len(args) > 0 && len(args[0]) > 1 {
			targetPath += args[0][1:]
		}
	case args[0] == "-":
		State.ToPreDir()
		return
	case strings.HasPrefix(args[0], "/"):
		targetPath = args[0]
	default:
		targetPath = path.Join(State.Cwd(), args[0])
	}
	if !IsExist(targetPath) {
		writeError(fmt.Sprintf("cd: %s: No such file or directory", args[0]))
		return
	}

	State.Chdir(targetPath)
}

func Echo(args []string) {
	writeOutput(strings.Join(args, " "))
}

func Type(args []string) {
	bin := args[0]
	if IsBuiltin(bin) {
		writeOutput(bin + " is a shell builtin")
		return
	}

	if file, exists := findBinInPath(bin); exists {
		writeOutput(fmt.Sprintf("%s is %s", bin, file))
		return
	}

	writeError(bin + ": not found")
}

func RunExternalApp(executable string, args []string) {
	_, found := findBinInPath(executable)
	if found {
		cmd := exec.Command(executable, args...)
		cmd.Stdout = State.GetOutputWriter()
		cmd.Stderr = State.GetErrorWriter()
		cmd.Run()
	} else {
		writeError(executable + ": command not found")
	}
}
