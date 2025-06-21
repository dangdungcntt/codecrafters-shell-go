package main

import (
	"bufio"
	"fmt"
	"github.com/codecrafters-io/shell-starter-go/internal/commands"
	"os"
	"strings"
	"unicode"
)

func main() {
	commands.Init()

	for {
		fmt.Fprint(os.Stdout, "$ ")

		// Wait for user input
		raw, err := bufio.NewReader(os.Stdin).ReadString('\n')
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error reading input:", err)
			os.Exit(1)
		}

		argv := parseCommand(raw)

		var args []string
		if len(argv) > 1 {
			args = argv[1:]
		}

		commands.ExecuteCommand(argv[0], args)
	}
}

func parseCommand(raw string) []string {
	raw = strings.TrimSuffix(raw, "\n")

	var current strings.Builder
	var args []string
	var inSingleQuote, inDoubleQuote, escaped bool

	for _, c := range raw {

		switch {
		case escaped:
			current.WriteRune(c)
			escaped = false

		case c == '\\' && !inDoubleQuote && !inSingleQuote:
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
