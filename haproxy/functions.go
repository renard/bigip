package haproxy

import (
	"strings"
)

func comment(str string, indent int) string {
	idstr := strings.Repeat(" ", indent)
	lines := strings.Split(str, "\n")
	for i, line := range lines {
		lines[i] = idstr + "# " + line
	}
	return strings.Join(lines, "\n")
}

func ipport(str string) string {
	strs := strings.Split(str, "/")
	return strs[len(strs)-1]
}
