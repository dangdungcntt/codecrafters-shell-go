package commands

import (
	"io"
	"log"
	"os"
)

type ShellState struct {
	dirHistories  []string
	cwd           string
	currentOutput io.Writer
}

func NewShellState() *ShellState {
	cwd, _ := os.Getwd()
	return &ShellState{
		cwd:           cwd,
		currentOutput: os.Stdout,
	}
}

func (s *ShellState) GetOutput() io.Writer {
	return s.currentOutput
}

func (s *ShellState) SetOutput(w io.Writer) io.Writer {
	c := s.currentOutput
	s.currentOutput = w
	return c
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

func ExecuteCommand(executable string, args []string, output io.Writer) {
	if output != nil {
		old := State.SetOutput(output)
		defer func() {
			State.SetOutput(old)
		}()
	}
	executor, found := CommandMap[executable]
	if !found {
		RunExternalApp(executable, args, State.GetOutput())
		return
	}

	executor(args)
}
