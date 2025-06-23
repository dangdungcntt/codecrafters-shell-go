package commands

import (
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
	if len(args) == 0 {
		histories = CommandHistory
	} else {
		limit, _ := strconv.Atoi(args[0])
		baseIndex = limit + 1
		histories = CommandHistory[max(0, len(CommandHistory)-limit):]
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
