package commands

import (
	"bufio"
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

	initHistory()

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

func initHistory() {
	// Get history file path
	filePath := os.Getenv("HISTFILE")
	if filePath == "" {
		return
	}
	nFile, err := os.OpenFile(filePath, os.O_RDONLY, 0644)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}

	scanner := bufio.NewScanner(nFile)
	for scanner.Scan() {
		AddCommandToHistory(scanner.Text())
	}

	LastInitHistoryIndex = len(CommandHistory) - 1
}

func WriteHistory(isAppend bool) {
	// Get history file path
	filePath := os.Getenv("HISTFILE")
	if filePath == "" {
		return
	}
	flags := os.O_RDWR | os.O_CREATE
	if isAppend {
		flags = flags | os.O_APPEND
	} else {
		flags = flags | os.O_TRUNC
	}
	nFile, err := os.OpenFile(filePath, flags, 0644)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer nFile.Close()

	var message string
	for i, line := range CommandHistory {
		if isAppend && i <= LastInitHistoryIndex {
			continue
		}
		if i != len(CommandHistory)-1 {
			message += line + "\n"
		} else {
			message += line
		}
	}
	_, err = fmt.Fprintln(nFile, message)
	if err != nil {
		fmt.Println("Error write file:", err)
	}
}

type CommandHandler interface {
	Execute(ctx *ShellContext, args []string)
}

var CommandMap = map[string]CommandHandler{}
var CommandHistory = make([]string, 0, 10)
var LastAppendHistoryIndex = -1
var LastInitHistoryIndex = -1

func AddCommandToHistory(cmd string) {
	CommandHistory = append(CommandHistory, cmd)
}

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
