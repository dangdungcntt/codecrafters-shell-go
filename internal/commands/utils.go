package commands

import (
	"fmt"
	"os"
	"strings"
)

func findBinInPath(bin string) (string, bool) {
	paths := os.Getenv("PATH")
	for _, path := range strings.Split(paths, ":") {
		file := path + "/" + bin
		if IsExist(file) {
			return file, true
		}
	}

	return "", false
}

func IsExist(file string) bool {
	_, err := os.Stat(file)
	return err == nil
}

func assertNoError(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func writeToConsole(args ...string) {
	for _, arg := range args {
		fmt.Fprint(os.Stdout, arg)
	}

	fmt.Fprint(os.Stdout, "\n")
}
