package commands

import (
	"fmt"
	"log"
	"os"
)

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

func assertNoError(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func writeToConsole(args ...string) {
	for _, arg := range args {
		fmt.Fprint(os.Stdout, arg)
	}

	fmt.Fprint(os.Stdout, "\n")
}
