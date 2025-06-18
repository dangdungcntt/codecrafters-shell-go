package main

import (
	"bufio"
	"fmt"
	"github.com/codecrafters-io/shell-starter-go/internal/commands"
	"os"
)

func main() {
	for {
		fmt.Fprint(os.Stdout, "$ ")

		// Wait for user input
		raw, err := bufio.NewReader(os.Stdin).ReadString('\n')
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error reading input:", err)
			os.Exit(1)
		}

		cmd := commands.NewCommand(raw)
		cmd.Execute()
	}
}
