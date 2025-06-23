package commands

import (
	"fmt"
	"github.com/chzyer/readline"
)

const TerminalBell = "\x07"

type CustomCompleter struct {
	PrefixCompleter *readline.PrefixCompleter
	lastPrefix      string
	tabCount        int
}

func (c *CustomCompleter) Do(line []rune, pos int) ([][]rune, int) {
	newLine, _ := c.PrefixCompleter.Do(line, pos)
	matches := len(newLine)
	word := string(line[:pos])
	if matches == 0 {
		fmt.Print(TerminalBell)
		return nil, 0
	}
	if matches == 1 {
		c.tabCount = 0
		return nil, 0
	}

	longestPrefix := longestCommonPrefix(newLine)
	if longestPrefix != "" {
		return [][]rune{
			[]rune(longestPrefix),
		}, 0
	}

	if word == c.lastPrefix {
		c.tabCount++
	} else {
		c.tabCount = 1
	}

	c.lastPrefix = word
	if c.tabCount == 1 {
		fmt.Print(TerminalBell)
		return nil, 0
	}

	fmt.Println()

	for i, runes := range newLine {
		fmt.Printf("%s%s", word, string(runes))
		if i < matches-1 {
			fmt.Printf(" ")
		}
	}

	fmt.Println()

	fmt.Print("$ " + word)
	c.tabCount = 0

	return nil, 0
}

func longestCommonPrefix(strs [][]rune) string {
	if len(strs) == 0 {
		return ""
	}

	prefix := string(strs[0])
	for _, s := range strs[1:] {
		for len(prefix) > 0 && len(s) < len(prefix) || string(s[:len(prefix)]) != prefix {
			prefix = prefix[:len(prefix)-1]
		}
		if prefix == "" {
			break
		}
	}
	return prefix
}
