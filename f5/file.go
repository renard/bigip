// Copyright © 2023 Sébastien Gross <seb•ɑƬ•chezwam•ɖɵʈ•org>
//
// Created: 2021-12-19
// Last changed: 2024-10-09 01:27:38
//
// This program is free software: you can redistribute it and/or
// modify it under the terms of the GNU Affero General Public License
// as published by the Free Software Foundation, either version 3 of
// the License, or (at your option) any later version.
//
// This program is distributed in the hope that it will be useful, but
// WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU
// Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public
// License along with this program. If not, see
// <http://www.gnu.org/licenses/>.

package f5

import (
	"fmt"
	"io/ioutil"
	"reflect"
	"strings"

	"bigip/internal/log"

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
	Lines [3]int
	// File where block is defined
	File string
}

// MatchPrefix checks that Content field of ParsedConfig matches at
// least one of the give prefixes.
func (c *ParsedConfig) MatchPrefix(prefixes ...string) bool {
	for _, prefix := range prefixes {
		if strings.HasPrefix(c.Content, prefix) {
			return true
		}
	}
	return false
}

type f5config map[string]F5Object
type F5Config struct {
	LtmNode        f5config
	LtmPool        f5config
	LtmVirtual     f5config
	LtmRule        f5config
	LtmPolicy      f5config
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
		LtmPolicy:      f5config{},
		LtmProfile:     f5config{},
		LtmMonitor:     f5config{},
		LtmPersistence: f5config{},
	}
}

// countBraces counts the braces balance in a line.
//
// If return value is:
//
//	0 the braces are balanced in the line.
//	>0 there are more opening braces than closing
//	<0 there are less opening braces than closing
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

func parseLines(l *log.Log, file, content string) (pc []ParsedConfig, err error) {
	lines := strings.Split(content, "\n")
	l.Debug("Need to process %d lines", len(lines))
	tmp := []string{}
	processed := 0
	opened := 0
	retCur := 0
	for i, line := range lines {
		if strings.HasPrefix(line, "ltm ") {
			l.Trace("%.8d %.4d %s", i, opened, line)
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
				// i is the line number. Since loop index starts at 0 and line
				// number starts at 1, we should add 1 to loop index.
				//
				// Last line is then i + 1.
				//
				// l is the block lines number. Since it counts the actual
				// line numbers we need to substract l - 1 from the last line
				// to reach the first line of the block.
				//
				// The first line number is then: ( i + 1 ) - ( l - 1 ), which
				// is in fact: i - l + 2
				//
				//    line  | i    | l  |
				//    100   | 99   | 1  | This is a
				//    101   | 100  | 2  | multi-line
				//    102   | 101  | 3  | match
				//
				// 102 = 101 + 1 (last line)
				// 100 = 101 + 1 - ( 3 - 1)  = 101 - 3 + 2 = 100 (first line)
				Lines: [3]int{i - l + 2, i + 1, l},
			})
			retCur += 1
			tmp = tmp[:0]
		case opened < 0:
			l.Error("line %d: opened braces: %d", i, opened)
			return
		case opened > 0:
			processed += 1
			tmp = append(tmp, line)
		}
	}
	l.Debug("Found %d blocks in %d lines", len(pc), processed)
	return

}

func ParseFile(l *log.Log, files []string) (cfg F5Config, err error) {
	cfg = NewF5Config()
	for _, file := range files {
		l.Debug("Parsing %s", file)
		tmpCfg, e := parseFile(l, file)
		if e != nil {
			l.Error("Error while parsing file %s: %s", file, e)
			return NewF5Config(), e
		}
		err = cfg.Merge(tmpCfg)
		if err != nil {
			l.Error("Error while merging %s: %s", file, e)
			return NewF5Config(), err
		}
	}
	return
}

// ParseFile read and split file.
func parseFile(l *log.Log, file string) (cfg F5Config, err error) {

	content, err := ioutil.ReadFile(file)
	if err != nil {
		return
	}

	pc, err := parseLines(l, file, string(content))

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
		case o.MatchPrefix("ltm node ", "node "):
			obj, e = newLtmNode(o)
			dest = cfg.LtmNode
		case o.MatchPrefix("ltm pool ", "pool "):
			obj, e = newLtmPool(o)
			dest = cfg.LtmPool
		case o.MatchPrefix("ltm virtual ", "virtual "):
			obj, e = newLtmVirtual(o)
			dest = cfg.LtmVirtual
		case o.MatchPrefix("ltm rule ", "rule "):
			obj, e = newLtmRule(o)
			dest = cfg.LtmRule
		case o.MatchPrefix("ltm policy ", "policy "):
			obj, e = newLtmPolicy(o)
			dest = cfg.LtmPolicy
		case o.MatchPrefix("ltm profile ", "profile "):
			obj, e = newLtmProfile(o)
			dest = cfg.LtmProfile
		case o.MatchPrefix("ltm monitor ", "monitor "):
			obj, e = newLtmMonitor(o)
			dest = cfg.LtmMonitor
		case o.MatchPrefix("ltm persistence ", "persistence "):
			obj, e = newLtmPersistence(o)
			dest = cfg.LtmPersistence
		}
		if e != nil {
			l.Error("%s: %s", strings.Split(o.Content, "\n")[0], e)
			continue
		}
		if obj != nil {
			dest[obj.GetName()] = obj
			obj = nil
		}
	}

	l.Info("Parsed %d objects %d lines: %d nodes, %d pools, %d virtuals, %d rules, %d profiles",
		len(pc), lines, len(cfg.LtmNode),
		len(cfg.LtmPool), len(cfg.LtmVirtual),
		len(cfg.LtmRule), len(cfg.LtmProfile),
	)

	if false {
		repr.Println(cfg)
	}
	return
}
