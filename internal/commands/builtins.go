package commands

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path"
	"strconv"
	"strings"
)

type BuiltinHandler func(ctx *ShellContext, args []string)

func (b BuiltinHandler) Execute(ctx *ShellContext, args []string) {
	b(ctx, args)
}

func Exit(ctx *ShellContext, args []string) {
	code, err := strconv.Atoi(args[0])
	ctx.AssertNoError(err)
	os.Exit(code)
}

func Pwd(ctx *ShellContext, _ []string) {
	ctx.WriteError(ctx.Cwd())
}

func Cd(ctx *ShellContext, args []string) {
	var targetPath string
	switch {
	case len(args) == 0 || strings.HasPrefix(args[0], "~"):
		targetPath, _ = os.UserHomeDir()
		if len(args) > 0 && len(args[0]) > 1 {
			targetPath += args[0][1:]
		}
	case args[0] == "-":
		ctx.ToPreDir()
		return
	case strings.HasPrefix(args[0], "/"):
		targetPath = args[0]
	default:
		targetPath = path.Join(ctx.Cwd(), args[0])
	}
	if !IsExist(targetPath) {
		ctx.WriteError(fmt.Sprintf("cd: %s: No such file or directory", args[0]))
		return
	}

	ctx.Chdir(targetPath)
}

func Echo(ctx *ShellContext, args []string) {
	ctx.WriteOutput(strings.Join(args, " "))
}

func Type(ctx *ShellContext, args []string) {
	bin := args[0]
	if IsBuiltin(bin) {
		ctx.WriteOutput(bin + " is a shell builtin")
		return
	}

	if file, exists := findBinInPath(bin); exists {
		ctx.WriteOutput(fmt.Sprintf("%s is %s", bin, file))
		return
	}

	ctx.WriteError(bin + ": not found")
}

func History(ctx *ShellContext, args []string) {
	var histories []string
	var baseIndex int
	switch {
	case len(args) == 0:
		histories = CommandHistory
	case len(args) == 1:
		limit, _ := strconv.Atoi(args[0])
		baseIndex = limit + 1
		histories = CommandHistory[max(0, len(CommandHistory)-limit):]
	case len(args) == 2 && args[0] == "-r":
		// read history from file
		file, err := os.OpenFile(args[1], os.O_APPEND|os.O_CREATE|os.O_RDONLY, 0644)
		if err != nil {
			fmt.Println("Error opening file:", err)
			return
		}

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			AddCommandToHistory(scanner.Text())
		}
	case len(args) == 2 && (args[0] == "-w" || args[0] == "-a"):
		flags := os.O_WRONLY | os.O_CREATE
		isAppend := false
		if args[0] == "-a" {
			isAppend = true
			flags = flags | os.O_APPEND
		} else {
			flags = flags | os.O_TRUNC
		}
		file, err := os.OpenFile(args[1], flags, 0644)
		if err != nil {
			fmt.Println("Error opening file:", err)
			return
		}
		defer file.Close()

		var message string
		for i, line := range CommandHistory {
			if isAppend && i <= LastAppendHistoryIndex {
				continue
			}
			if i != len(CommandHistory)-1 {
				message += line + "\n"
			} else {
				message += line
			}
		}
		_, err = fmt.Fprintln(file, message)
		if err != nil {
			fmt.Println("Error writing to file:", err)
			return
		}
		if isAppend {
			LastAppendHistoryIndex = len(CommandHistory) - 1
		}
	}
	for i, cmd := range histories {
		ctx.WriteOutput(fmt.Sprintf("%5d %s", baseIndex+i+1, cmd))
	}
}

func RunExternalApp(ctx *ShellContext, executable string, args []string) {
	_, found := findBinInPath(executable)
	if found {
		cmd := exec.Command(executable, args...)
		cmd.Stdin = ctx.stdin
		cmd.Stdout = ctx.stdout
		cmd.Stderr = ctx.stderr
		cmd.Run()
	} else {
		ctx.WriteError(executable + ": command not found")
	}
}
