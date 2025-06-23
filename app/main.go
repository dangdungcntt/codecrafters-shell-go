package main

import (
	"cmp"
	"fmt"
	"github.com/chzyer/readline"
	"github.com/codecrafters-io/shell-starter-go/internal/commands"
	"io"
	"log"
	"os"
	"slices"
	"strings"
	"sync"
	"unicode"
)

func main() {
	allCommands := commands.Init()
	slices.SortStableFunc(allCommands, func(a, b string) int {
		return cmp.Compare(a, b)
	})

	completerList := make([]readline.PrefixCompleterInterface, 0, len(allCommands)+1)
	for _, cmd := range allCommands {
		completerList = append(completerList, readline.PcItem(cmd))
	}

	autoCompleter := &commands.CustomCompleter{
		PrefixCompleter: readline.NewPrefixCompleter(completerList...),
	}
	l, err := readline.NewEx(&readline.Config{
		Prompt:       "$ ",
		AutoComplete: autoCompleter,
	})
	if err != nil {
		log.Fatal(err)
	}

	for {
		fmt.Fprint(os.Stdout, "$ ")

		// Wait for user input
		raw, err := l.Readline()
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error reading input:", err)
			os.Exit(1)
		}
		commands.AddCommandToHistory(raw)

		argv := parseCommand(raw)

		cmds, err := parsePipeCommand(argv)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error parsePipeCommand:", err)
			os.Exit(1)
		}

		if len(cmds) == 1 {
			outputFile, errorFile := cmds[0].stdout, cmds[0].stderr
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error creating file: %v\n", err)
				os.Exit(1)
			}
			commands.ExecuteCommand(cmds[0].executable, cmds[0].args, nil, outputFile, errorFile)
		} else {
			wg := sync.WaitGroup{}
			wg.Add(len(cmds))
			var lastStdin io.ReadCloser
			createdFiles := make([]io.Closer, 0, (len(cmds)-1)*2)
			for i := 0; i < len(cmds); i++ {
				cmd := cmds[i]

				var cmdOutput io.WriteCloser
				cmdStdin := lastStdin
				if i == len(cmds)-1 {
					cmdOutput = cmd.stdout
				} else {
					newR, newW, _ := os.Pipe()
					createdFiles = append(createdFiles, newR)
					createdFiles = append(createdFiles, newW)
					lastStdin = newR
					cmdOutput = newW
				}
				go func(stdin io.Reader, stdout io.WriteCloser, stderr io.Writer) {
					defer func() {
						wg.Done()
						if stdout != nil {
							stdout.Close()
						}
					}()
					commands.ExecuteCommand(cmd.executable, cmd.args, stdin, stdout, stderr)
				}(cmdStdin, cmdOutput, cmd.stderr)
			}

			wg.Wait()

			for _, file := range createdFiles {
				file.Close()
			}
		}
	}
}

type SingleCommand struct {
	executable string
	args       []string
	stdout     io.WriteCloser
	stderr     io.Writer
}

func NewSingleCommand(executable string, args []string) (*SingleCommand, error) {
	var outputFile, errorFile io.WriteCloser
	var err error
	for i, arg := range args {
		if (arg == ">" || arg == "1>") && i+1 < len(args) {
			if outputFile, err = os.Create(args[i+1]); err != nil {
				return nil, err
			}
			args = args[:i]
			break
		}

		if (arg == ">>" || arg == "1>>") && i+1 < len(args) {
			if outputFile, err = os.OpenFile(args[i+1], os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666); err != nil {
				return nil, err
			}
			args = args[:i]
			break
		}

		if arg == "2>" && i+1 < len(args) {
			if errorFile, err = os.Create(args[i+1]); err != nil {
				return nil, err
			}
			args = args[:i]
			break
		}

		if arg == "2>>" && i+1 < len(args) {
			if errorFile, err = os.OpenFile(args[i+1], os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666); err != nil {
				return nil, err
			}
			args = args[:i]
			break
		}
	}

	return &SingleCommand{
		executable: executable,
		args:       args,
		stdout:     outputFile,
		stderr:     errorFile,
	}, nil
}

func parsePipeCommand(argv []string) ([]*SingleCommand, error) {
	lastStartIndex := 0
	results := make([]*SingleCommand, 0)
	for i, arg := range argv {
		if arg == "|" {
			args := argv[lastStartIndex:i]
			cmd, err := NewSingleCommand(args[0], nil)
			if err != nil {
				return nil, err
			}
			if len(args) > 1 {
				cmd.args = args[1:]
			}
			results = append(results, cmd)
			lastStartIndex = i + 1
		}
	}

	if lastStartIndex < len(argv) {
		argv = argv[lastStartIndex:]
		exe := argv[0]
		var args []string
		if len(argv) > 1 {
			args = argv[1:]
		}
		cmd, err := NewSingleCommand(exe, args)
		if err != nil {
			return nil, err
		}

		results = append(results, cmd)
	}

	return results, nil
}

func parseCommand(raw string) []string {
	raw = strings.TrimSuffix(raw, "\n")

	var current strings.Builder
	var args []string
	var inSingleQuote, inDoubleQuote, escaped bool

	for _, c := range raw {

		switch {
		case escaped && inDoubleQuote:
			if c == '$' || c == '"' || c == '\\' {
				current.WriteRune(c)
			} else {
				current.WriteRune('\\')
				current.WriteRune(c)
			}
			escaped = false
		case escaped:
			current.WriteRune(c)
			escaped = false

		case c == '\\' && !inSingleQuote:
			escaped = true

		case c == '\'' && !inDoubleQuote:
			inSingleQuote = !inSingleQuote

		case c == '"' && !inSingleQuote:
			inDoubleQuote = !inDoubleQuote

		case unicode.IsSpace(c) && !inSingleQuote && !inDoubleQuote:
			if current.Len() > 0 {
				args = append(args, current.String())
				current.Reset()
			}

		default:
			current.WriteRune(c)

		}
	}

	if current.Len() > 0 {
		args = append(args, current.String())
	}

	return args
}
