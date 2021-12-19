package f5

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/alecthomas/repr"
)

type F5Object interface {
	Original() string
}

type f5config map[string]F5Object
type F5Config struct {
	LtmNode    f5config
	LtmPool    f5config
	LtmVirtual f5config
}

func newConfigOject(data []string) (ret F5Config) {
	// kw := strings.Fields(data[0])
	// fmt.Printf("kw: %#v %d\n", kw, len(kw))
	// switch {
	// case len(kw) < 2:
	// 	return
	// case kw[0] == "ltm":
	// 	switch {
	// 	case kw[1] == "virtual":
	// 		ret = newLtmVirtual(data)
	// 	}
	// }
	return ret

}

// countBraces counts the braces balance in a line.
//
// If return value is:
//  0 the braces are balanced in the line.
//  >0 there are more opening bbraces than closing
//  <0 there are less opening bbraces than closing
func countBraces(line string) (count int) {
	count = 0
	for _, c := range line {
		switch c {
		// case '\'':
		// 	simple_quote = !simple_quote
		// case '"':
		// 	double_quote = !double_quote
		// // Do not include braces in comments
		// case '#':
		// 	return
		case '{':
			count += 1
		case '}':
			count -= 1
		}
	}
	return
}

func countBraces2(line string) (count int) {
	count = 0
	simple_quote := false
	double_quote := false
	for _, c := range line {
		switch c {
		case '\'':
			simple_quote = !simple_quote
		case '"':
			double_quote = !double_quote
		// Do not include braces in comments
		case '#':
			return
		case '{':
			if !simple_quote && !double_quote {
				count += 1
			}
		case '}':
			if !simple_quote && !double_quote {
				count -= 1
			}
		}
	}
	return
}

type ParsedConfig struct {
	Content []string
	Lines   [2]int
	File    string
}

func parseLines(content string) (pc []ParsedConfig, err error) {
	lines := strings.Split(content, "\n")
	fmt.Printf("Need to process %d lines\n", len(lines))
	tmp := []string{}
	processed := 0
	opened := 0
	retCur := 0
	for i, line := range lines {
		// fmt.Printf("%d\t%s\n", opened, line)
		if false {
			if strings.HasPrefix(line, "ltm ") && opened != 0 {
				// fmt.Printf("Unbalanced braces line %d (offset: %d)\n", i, opened)
				opened = 0
			} else {
				opened += countBraces(line)
			}
		}
		opened += countBraces(line)
		switch {
		case opened == 0:
			switch {
			case line == "":
				continue
			case line[0] == '#':
				continue
			}
			tmp = append(tmp, line)
			l := len(tmp)
			processed += 1
			pc = append(pc, ParsedConfig{
				Content: make([]string, l),
				Lines:   [2]int{i - l, i},
			})
			copy(pc[retCur].Content, tmp)

			retCur += 1
			tmp = tmp[:0]
		case opened < 0:
			fmt.Printf("Error line %d: opened braces: %d\n", i, opened)
			return
		case opened > 0:
			processed += 1
			tmp = append(tmp, line)
		}
	}

	// for _, o := range ret {
	// 	fmt.Println(o[0])
	// 	t := newConfigOject(o)
	// 	fmt.Printf("%#v\n", t)
	// }
	fmt.Printf("Found %d blocks in %d lines\n", len(pc), processed)
	return

}

// ParseFile read and split file.
func ParseFile(file string) (cfg F5Config, err error) {

	content, err := ioutil.ReadFile(file)
	if err != nil {
		return
	}

	pc, err := parseLines(string(content))

	cfg = F5Config{
		LtmNode:    f5config{},
		LtmVirtual: f5config{},
		LtmPool:    f5config{},
	}

	nodes := 0
	pools := 0
	virtuals := 0
	lines := 0
	for _, o := range pc {
		lines += len(o.Content)
		// fmt.Printf("%d\t%s\n", len(o), o[0])
		switch {
		case strings.HasPrefix(o.Content[0], "ltm node "):
			obj, e := newLtmNode(strings.Join(o.Content, "\n"))
			if e != nil {
				fmt.Printf("Err: %s: %s\n", o.Content[0], e)
				continue
			}
			cfg.LtmNode[obj.Name] = obj
			_ = obj
			nodes += 1
		case strings.HasPrefix(o.Content[0], "ltm pool "):
			obj, e := newLtmPool(strings.Join(o.Content, "\n"))
			if e != nil {
				fmt.Printf("Err: %s: %s\n", o.Content[0], e)
				continue
			}
			cfg.LtmPool[obj.Name] = obj
			_ = obj
			pools += 1
		case strings.HasPrefix(o.Content[0], "ltm virtual "):
			obj, e := newLtmVirtual(strings.Join(o.Content, "\n"))
			if e != nil {
				fmt.Printf("Err: %s: %s\n", o.Content[0], e)
				continue
			}
			cfg.LtmVirtual[obj.Name] = obj
			_ = obj
			virtuals += 1
		}
	}

	fmt.Printf("Parsed %d objects %d lines: %d nodes, %d pools, %d virtuals\n", len(pc), lines, nodes, pools, virtuals)
	if false {
		repr.Println(cfg)
	}
	return
}
