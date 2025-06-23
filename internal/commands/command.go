package commands

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

type ShellState struct {
	dirHistories []string
	cwd          string
	outputWriter io.Writer
	errorWriter  io.Writer
}

func NewShellState() *ShellState {
	cwd, _ := os.Getwd()
	return &ShellState{
		cwd:          cwd,
		outputWriter: os.Stdout,
		errorWriter:  os.Stderr,
	}
}

func (s *ShellState) GetOutputWriter() io.Writer {
	return s.outputWriter
}

func (s *ShellState) SetOutputWriter(w io.Writer) io.Writer {
	c := s.outputWriter
	s.outputWriter = w
	return c
}

func (s *ShellState) GetErrorWriter() io.Writer {
	return s.errorWriter
}

func (s *ShellState) SetErrorWriter(w io.Writer) io.Writer {
	c := s.errorWriter
	s.errorWriter = w
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

func Init() []string {
	RegisterCommand("echo", Echo)
	RegisterCommand("type", Type)
	RegisterCommand("exit", Exit)
	RegisterCommand("pwd", Pwd)
	RegisterCommand("cd", Cd)

	allCommandsMap := make(map[string]struct{}, len(CommandMap))
	for c := range CommandMap {
		allCommandsMap[c] = struct{}{}
	}

	for _, dirPath := range strings.Split(os.Getenv("PATH"), ":") {
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
		}
	}

	allCommands := make([]string, 0, len(allCommandsMap))
	for c := range allCommandsMap {
		allCommands = append(allCommands, c)
	}

	return allCommands
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

func ExecuteCommand(executable string, args []string, outputWriter io.Writer, errorWriter io.Writer) {
	if outputWriter != nil {
		old := State.SetOutputWriter(outputWriter)
		defer func() {
			State.SetOutputWriter(old)
		}()
	}

	if errorWriter != nil {
		old := State.SetErrorWriter(errorWriter)
		defer func() {
			State.SetErrorWriter(old)
		}()
	}

	executor, found := CommandMap[executable]
	if !found {
		RunExternalApp(executable, args)
		return
	}

	executor(args)
}
