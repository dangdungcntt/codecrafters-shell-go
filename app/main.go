package main

import (
	"bufio"
	"fmt"
	"github.com/codecrafters-io/shell-starter-go/internal/commands"
	"os"
	"strings"
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

func parseCommand(cmd string) []string {
	cmd = strings.TrimSuffix(cmd, "\n")

	var args []string
	var currentArg strings.Builder
	var inSingleQuote, inDoubleQuote bool
	var escapeNext bool

	for _, char := range cmd {
		if escapeNext {
			currentArg.WriteRune(char)
			escapeNext = false
			continue
		}

		switch char {
		case '\\':
			escapeNext = true
		case '\'':
			if !inDoubleQuote {
				inSingleQuote = !inSingleQuote
			} else {
				currentArg.WriteRune(char)
			}
		case '"':
			if !inSingleQuote {
				inDoubleQuote = !inDoubleQuote
			} else {
				currentArg.WriteRune(char)
			}
		case ' ':
			if inSingleQuote || inDoubleQuote {
				currentArg.WriteRune(char)
			} else {
				if currentArg.Len() > 0 {
					args = append(args, currentArg.String())
					currentArg.Reset()
				}
			}
		default:
			currentArg.WriteRune(char)
		}
	}

	if currentArg.Len() > 0 {
		args = append(args, currentArg.String())
	}

	return args
}
