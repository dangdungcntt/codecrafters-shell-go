package commands

import (
	"fmt"
	"io"
	"os"
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
