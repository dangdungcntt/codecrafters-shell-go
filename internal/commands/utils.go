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

func writeOutput(args ...string) {
	outputWriter := State.GetOutputWriter()
	for _, arg := range args {
		fmt.Fprint(outputWriter, arg)
	}

	fmt.Fprint(outputWriter, "\n")
}

func writeError(args ...string) {
	errorWriter := State.GetErrorWriter()
	for _, arg := range args {
		fmt.Fprint(errorWriter, arg)
	}

	fmt.Fprint(errorWriter, "\n")
}
