package commands

import (
	"log"
	"os"
)

var State = struct {
	Pwd string
}{}

func Init() {
	RegisterCommand("echo", Echo)
	RegisterCommand("type", Type)
	RegisterCommand("exit", Exit)
	RegisterCommand("pwd", Pwd)
	RegisterCommand("cd", Cd)
	State.Pwd, _ = os.Getwd()
}

var CommandMap = map[string]func(args []string){}

func IsBuiltin(name string) bool {
	_, found := CommandMap[name]
	return found
}

func RegisterCommand(name string, executor func(args []string)) {
	_, found := CommandMap[name]
	if found {
		log.Fatal("command " + name + " already existed")
	}

	CommandMap[name] = executor
}

func ExecuteCommand(executable string, args []string) {
	executor, found := CommandMap[executable]
	if !found {
		RunExternalApp(executable, args)
		return
	}

	executor(args)
}
