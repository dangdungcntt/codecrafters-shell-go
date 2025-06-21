package commands

import (
	"log"
	"os"
)

type ShellState struct {
	dirHistories []string
	cwd          string
}

func NewShellState() *ShellState {
	cwd, _ := os.Getwd()
	return &ShellState{
		cwd: cwd,
	}
}

func (s *ShellState) Chdir(path string) {
	assertNoError(os.Chdir(path))
	s.dirHistories = append(s.dirHistories, s.cwd)
	s.cwd = path
}

func (s *ShellState) ToPreDir() {
	if len(s.dirHistories) == 0 {
		return
	}

	s.Chdir(s.dirHistories[len(s.dirHistories)-1])
}

func (s *ShellState) Cwd() string {
	return s.cwd
}

var State = NewShellState()

func Init() {
	RegisterCommand("echo", Echo)
	RegisterCommand("type", Type)
	RegisterCommand("exit", Exit)
	RegisterCommand("pwd", Pwd)
	RegisterCommand("cd", Cd)
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
