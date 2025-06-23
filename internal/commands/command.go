package commands

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

func Init() []string {
	RegisterCommand("echo", BuiltinHandler(Echo))
	RegisterCommand("type", BuiltinHandler(Type))
	RegisterCommand("exit", BuiltinHandler(Exit))
	RegisterCommand("pwd", BuiltinHandler(Pwd))
	RegisterCommand("cd", BuiltinHandler(Cd))
	RegisterCommand("history", BuiltinHandler(History))

	allCommandsMap := make(map[string]struct{}, len(CommandMap))
	allCommands := make([]string, 0, len(allCommandsMap))
	for c := range CommandMap {
		allCommandsMap[c] = struct{}{}
		allCommands = append(allCommands, c)
	}

	for _, dirPath := range strings.Split(os.Getenv("PATH"), string(os.PathListSeparator)) {
		files, err := os.ReadDir(dirPath)
		if err != nil {
			fmt.Errorf("error during reading dir: %v \n", err.Error())
			continue
		}

		for _, entry := range files {
			name := entry.Name()
			_, ok := allCommandsMap[name]
			if ok {
				continue
			}
			allCommandsMap[name] = struct{}{}
			allCommands = append(allCommands, name)
		}
	}

	return allCommands
}

type CommandHandler interface {
	Execute(ctx *ShellContext, args []string)
}

var CommandMap = map[string]CommandHandler{}

func IsBuiltin(name string) bool {
	_, found := CommandMap[name]
	return found
}

func RegisterCommand(name string, executor CommandHandler) {
	_, found := CommandMap[name]
	if found {
		log.Fatal("command " + name + " already existed")
	}

	CommandMap[name] = executor
}

func ExecuteCommand(executable string, args []string, stdin io.Reader, stdout io.Writer, stderr io.Writer) {
	ctx := NewShellContext(stdin, stdout, stderr)
	cmd, found := CommandMap[executable]
	if !found {
		RunExternalApp(ctx, executable, args)
		return
	}

	cmd.Execute(ctx, args)
}
