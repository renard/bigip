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

type ParsedConfig struct {
	Content string
	Lines   [2]int
	File    string
}

type f5config map[string]F5Object
type F5Config struct {
	LtmNode    f5config
	LtmPool    f5config
	LtmVirtual f5config
	LtmRule    f5config
	LtmProfile f5config
}

func newF5Config() F5Config {
	return F5Config{
		LtmNode:    f5config{},
		LtmVirtual: f5config{},
		LtmPool:    f5config{},
		LtmRule:    f5config{},
		LtmProfile: f5config{},
	}
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
				Content: strings.Join(tmp, "\n"),
				Lines:   [2]int{i - l, i},
			})

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

	cfg = newF5Config()

	lines := 0
	for _, o := range pc {
		lines += len(o.Content)
		switch {
		case strings.HasPrefix(o.Content, "ltm node ") || strings.HasPrefix(o.Content, "node "):
			obj, e := newLtmNode(o)
			if e != nil {
				fmt.Printf("Err: %s: %s\n", strings.Split(o.Content, "\n")[0], e)
				continue
			}
			cfg.LtmNode[obj.Name] = obj
		case strings.HasPrefix(o.Content, "ltm pool ") || strings.HasPrefix(o.Content, "pool "):
			obj, e := newLtmPool(o)
			if e != nil {
				fmt.Printf("Err: %s: %s\n", strings.Split(o.Content, "\n")[0], e)
				continue
			}
			cfg.LtmPool[obj.Name] = obj
		case strings.HasPrefix(o.Content, "ltm virtual ") || strings.HasPrefix(o.Content, "virtual "):
			obj, e := newLtmVirtual(o)
			if e != nil {
				fmt.Printf("Err: %s: %s\n", strings.Split(o.Content, "\n")[0], e)
				continue
			}
			cfg.LtmVirtual[obj.Name] = obj
		case strings.HasPrefix(o.Content, "ltm rule ") || strings.HasPrefix(o.Content, "rule "):
			obj, e := newLtmRule(o)
			if e != nil {
				fmt.Printf("Err: %s: %s\n", strings.Split(o.Content, "\n")[0], e)
				continue
			}
			cfg.LtmRule[obj.Name] = obj
		case strings.HasPrefix(o.Content, "ltm profile ") || strings.HasPrefix(o.Content, "profile "):
			obj, e := newLtmProfile(o)
			if e != nil {
				fmt.Printf("Err: %s: %s\n", strings.Split(o.Content, "\n")[0], e)
				continue
			}
			cfg.LtmProfile[obj.Name] = obj
		}
	}

	fmt.Printf("Parsed %d objects %d lines: %d nodes, %d pools, %d virtuals, %d rules, %d profiles\n",
		len(pc), lines, len(cfg.LtmNode),
		len(cfg.LtmPool), len(cfg.LtmVirtual),
		len(cfg.LtmRule), len(cfg.LtmProfile),
	)

	if false {
		repr.Println(cfg)
	}
	return
}
