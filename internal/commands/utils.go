package commands

import (
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
