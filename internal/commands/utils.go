package commands

import (
	"os"
	"strings"
)

func findBinInPath(bin string) (string, bool) {
	paths := os.Getenv("PATH")
	for _, path := range strings.Split(paths, ":") {
		file := path + "/" + bin
		if _, err := os.Stat(file); err == nil {
			return file, true
		}
	}

	return "", false
}
