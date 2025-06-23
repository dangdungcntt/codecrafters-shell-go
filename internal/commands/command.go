package commands

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

type ShellContext struct {
	dirHistories []string
	cwd          string
	stdin        io.Reader
	stdout       io.Writer
	stderr       io.Writer
}

func NewShellContext(stdin io.Reader, stdout io.Writer, stderr io.Writer) *ShellContext {
	cwd, _ := os.Getwd()

	if stdout == nil {
		stdout = os.Stdout
	}

	if stderr == nil {
		stderr = os.Stderr
	}

	return &ShellContext{
		cwd:    cwd,
		stdin:  stdin,
		stdout: stdout,
		stderr: stderr,
	}
}

func (s *ShellContext) GetOutputWriter() io.Writer {
	return s.stdout
}

func (s *ShellContext) SetOutputWriter(w io.Writer) io.Writer {
	c := s.stdout
	s.stdout = w
	return c
}

func (s *ShellContext) GetErrorWriter() io.Writer {
	return s.stderr
}

func (s *ShellContext) SetErrorWriter(w io.Writer) io.Writer {
	c := s.stderr
	s.stderr = w
	return c
}

func (s *ShellContext) Chdir(path string) {
	s.AssertNoError(os.Chdir(path))
	s.dirHistories = append(s.dirHistories, s.cwd)
	s.cwd = path
}

func (s *ShellContext) ToPreDir() {
	if len(s.dirHistories) == 0 {
		return
	}

	s.Chdir(s.dirHistories[len(s.dirHistories)-1])
}

func (s *ShellContext) Cwd() string {
	return s.cwd
}

func (s *ShellContext) WriteOutput(args ...string) {
	for _, arg := range args {
		fmt.Fprint(s.stdout, arg)
	}

	fmt.Fprint(s.stdout, "\n")
}

func (s *ShellContext) WriteError(args ...string) {
	for _, arg := range args {
		fmt.Fprint(s.stderr, arg)
	}

	fmt.Fprint(s.stderr, "\n")
}

func (s *ShellContext) AssertNoError(err error) {
	if err != nil {
		fmt.Fprintln(s.stderr, err)
		os.Exit(1)
	}
}

func Init() []string {
	RegisterCommand("echo", BuiltinHandler(Echo))
	RegisterCommand("type", BuiltinHandler(Type))
	RegisterCommand("exit", BuiltinHandler(Exit))
	RegisterCommand("pwd", BuiltinHandler(Pwd))
	RegisterCommand("cd", BuiltinHandler(Cd))

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
