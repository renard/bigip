package f5

import (
	"fmt"
	"io/ioutil"
	"reflect"
	"strings"

	"github.com/alecthomas/repr"
)

type F5Object interface {
	Original() string
	GetName() string
}

// ParsedConfig containts basic inforation of a parsed block
type ParsedConfig struct {
	// Original block content
	Content string
	// Begin and end lines of block in File
	Lines [2]int
	// File where block is defined
	File string
}

type f5config map[string]F5Object
type F5Config struct {
	LtmNode        f5config
	LtmPool        f5config
	LtmVirtual     f5config
	LtmRule        f5config
	LtmProfile     f5config
	LtmMonitor     f5config
	LtmPersistence f5config
}

type DuplicatedConfigEntry struct {
	obj   string
	entry string
}

func (e DuplicatedConfigEntry) Error() string {
	return fmt.Sprintf("Duplicated config Entry: %s", e.entry)
}

// Info returns the number of objects of all F5Config fields. This is
// mainly for debuging purposes.
func (c *F5Config) Info() string {
	refc := reflect.Indirect(reflect.ValueOf(c))
	ret := make([]string, refc.NumField())
	for i := 0; i < refc.NumField(); i++ {
		field := refc.Type().Field(i).Name
		mapc := refc.FieldByName(field).Interface().(f5config)
		ret[i] = fmt.Sprintf("%d %s", len(mapc), field)
	}
	return strings.Join(ret, ", ")
}

// Merge meges F5config n into c and returns DuplicatedConfigEntry if
// an entry from n is already defined in c.
//
// For each F5Config field, it loop over all entries from n and insert
// them into n ad a new entry.
func (c *F5Config) Merge(n F5Config) error {
	// Define reflection for both c and n
	refc := reflect.Indirect(reflect.ValueOf(c))
	refn := reflect.Indirect(reflect.ValueOf(n))
	// Loop over each F5Config struct fields
	for i := 0; i < refc.NumField(); i++ {
		field := refc.Type().Field(i).Name
		// Assign a variable (type f5config) to each structure field
		mapc := refc.FieldByName(field).Interface().(f5config)
		mapn := refn.FieldByName(field).Interface().(f5config)
		// Copy data from n to c
		for k, v := range mapn {
			if ok := mapc[k]; ok != nil {
				return DuplicatedConfigEntry{entry: k, obj: field}
			}
			mapc[k] = v
		}
	}

	return nil
}

func NewF5Config() F5Config {
	return F5Config{
		LtmNode:        f5config{},
		LtmVirtual:     f5config{},
		LtmPool:        f5config{},
		LtmRule:        f5config{},
		LtmProfile:     f5config{},
		LtmMonitor:     f5config{},
		LtmPersistence: f5config{},
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
	// simple_quote := false
	// double_quote := false
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

func parseLines(file, content string) (pc []ParsedConfig, err error) {
	lines := strings.Split(content, "\n")
	if false {
		fmt.Printf("Need to process %d lines\n", len(lines))
	}
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
				File:    file,
				// Why?
				Lines: [2]int{i - l + 2, i + 1},
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
	if false {
		fmt.Printf("Found %d blocks in %d lines\n", len(pc), processed)
	}
	return

}

func ParseFile(files []string) (cfg F5Config, err error) {
	cfg = NewF5Config()
	for _, file := range files {
		tmpCfg, e := parseFile(file)
		if e != nil {
			return NewF5Config(), e
		}
		err = cfg.Merge(tmpCfg)
		if err != nil {
			return NewF5Config(), err
		}
	}
	return
}

// ParseFile read and split file.
func parseFile(file string) (cfg F5Config, err error) {

	content, err := ioutil.ReadFile(file)
	if err != nil {
		return
	}

	pc, err := parseLines(file, string(content))

	cfg = NewF5Config()

	lines := 0
	for _, o := range pc {
		lines += len(o.Content)
		var (
			dest f5config
			obj  F5Object
			e    error
		)
		switch {
		case strings.HasPrefix(o.Content, "ltm node ") || strings.HasPrefix(o.Content, "node "):
			obj, e = newLtmNode(o)
			dest = cfg.LtmNode
		case strings.HasPrefix(o.Content, "ltm pool ") || strings.HasPrefix(o.Content, "pool "):
			obj, e = newLtmPool(o)
			dest = cfg.LtmPool
		case strings.HasPrefix(o.Content, "ltm virtual ") || strings.HasPrefix(o.Content, "virtual "):
			obj, e = newLtmVirtual(o)
			dest = cfg.LtmVirtual
		case strings.HasPrefix(o.Content, "ltm rule ") || strings.HasPrefix(o.Content, "rule "):
			obj, e = newLtmRule(o)
			dest = cfg.LtmRule
		case strings.HasPrefix(o.Content, "ltm profile ") || strings.HasPrefix(o.Content, "profile "):
			obj, e = newLtmProfile(o)
			dest = cfg.LtmProfile
		case strings.HasPrefix(o.Content, "ltm monitor ") || strings.HasPrefix(o.Content, "monitor "):
			obj, e = newLtmMonitor(o)
			dest = cfg.LtmMonitor
		case strings.HasPrefix(o.Content, "ltm persistence ") || strings.HasPrefix(o.Content, "persistence "):
			obj, e = newLtmPersistence(o)
			dest = cfg.LtmPersistence
		}
		if e != nil {
			fmt.Printf("Err: %s: %s\n", strings.Split(o.Content, "\n")[0], e)
			continue
		}
		if obj != nil {
			dest[obj.GetName()] = obj
			obj = nil
		}
	}

	if false {
		fmt.Printf("Parsed %d objects %d lines: %d nodes, %d pools, %d virtuals, %d rules, %d profiles\n",
			len(pc), lines, len(cfg.LtmNode),
			len(cfg.LtmPool), len(cfg.LtmVirtual),
			len(cfg.LtmRule), len(cfg.LtmProfile),
		)
	}
	if false {
		repr.Println(cfg)
	}
	return
}
