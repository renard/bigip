package haproxy

import (
	"strings"
)

func comment(str string, indent int) string {
	return spacedComment(str, 1, indent)
}

func spacedComment(str string, space, indent int) string {
	idstr := strings.Repeat(" ", indent)
	spstr := strings.Repeat(" ", space)
	lines := strings.Split(str, "\n")
	for i, line := range lines {
		lines[i] = idstr + "#" + spstr + line
	}
	return strings.Join(lines, "\n")
}

func ipport(str string) string {
	strs := strings.Split(str, "/")
	return strs[len(strs)-1]
}
