package commands

import "strings"

func Echo(args []string) {
	writeToConsole(strings.Join(args, " "))
}
