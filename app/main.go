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

		argv := parseCommand(raw)

		var outputFile, errorFile io.WriteCloser
		for i, arg := range argv {
			if (arg == ">" || arg == "1>") && i+1 < len(argv) {
				if outputFile, err = os.Create(argv[i+1]); err != nil {
					fmt.Fprintf(os.Stderr, "Error creating file: %v\n", err)
					os.Exit(1)
				}
				argv = argv[:i]
				break
			}

			if (arg == ">>" || arg == "1>>") && i+1 < len(argv) {
				if outputFile, err = os.OpenFile(argv[i+1], os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666); err != nil {
					fmt.Fprintf(os.Stderr, "Error creating file: %v\n", err)
					os.Exit(1)
				}
				argv = argv[:i]
				break
			}

			if arg == "2>" && i+1 < len(argv) {
				if errorFile, err = os.Create(argv[i+1]); err != nil {
					fmt.Fprintf(os.Stderr, "Error creating file: %v\n", err)
					os.Exit(1)
				}
				argv = argv[:i]
				break
			}

			if arg == "2>>" && i+1 < len(argv) {
				if errorFile, err = os.OpenFile(argv[i+1], os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666); err != nil {
					fmt.Fprintf(os.Stderr, "Error creating file: %v\n", err)
					os.Exit(1)
				}
				argv = argv[:i]
				break
			}
		}

		var args []string
		if len(argv) > 1 {
			args = argv[1:]
		}

		commands.ExecuteCommand(argv[0], args, outputFile, errorFile)

		if outputFile != nil {
			outputFile.Close()
		}

		if errorFile != nil {
			errorFile.Close()
		}
	}
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
