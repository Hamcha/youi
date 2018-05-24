package utils

import (
	"strings"
)

// IndentStrings adds indentation to all lines in a string
func IndentStrings(str string, indentLevel int) string {
	indentstr := strings.Repeat("  ", indentLevel)
	return indentstr + strings.TrimSpace(strings.Replace(str, "\n", "\n"+indentstr, -1))
}
