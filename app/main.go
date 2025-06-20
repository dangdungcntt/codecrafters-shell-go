package main

import (
	"bufio"
	"errors"
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

		argv := strings.Fields(raw)

		var args []string
		if len(argv) > 1 {
			args = argv[1:]
		}

		commands.ExecuteCommand(argv[0], args)
	}
}

func parse(raw string) (string, string, error) {
	raw = strings.TrimSuffix(raw, "\n")
	parts := strings.SplitN(raw, " ", 2)
	if len(parts) == 0 {
		return "", "", errors.New("empty command")
	}

	var args string
	if len(parts) == 2 {
		args = parts[1]
	}

	return parts[0], args, nil
}
