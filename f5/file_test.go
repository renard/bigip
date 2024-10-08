// Copyright © 2023 Sébastien Gross <seb•ɑƬ•chezwam•ɖɵʈ•org>
//
// Created: 2022-01-06
// Last changed: 2024-10-09 01:24:12
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
	"embed"
	"reflect"
	"testing"

	"bigip/internal/log"

	"github.com/alecthomas/repr"
)

var (
	//go:embed testdata/bigip-*.conf
	testsF5config embed.FS
	testLog       = log.New()
)

func TestMerge(t *testing.T) {
	c1 := NewF5Config()
	c1.LtmNode["Foo"] = &LtmNode{Name: "Foo"}
	c1.LtmNode["Bar"] = &LtmNode{Name: "Bar"}
	c2 := NewF5Config()
	c2.LtmRule["Foo"] = &LtmRule{Name: "Foo"}
	c2.LtmRule["Bar"] = &LtmRule{Name: "Bar"}

	c3 := NewF5Config()
	c3.LtmNode["Foo"] = &LtmNode{Name: "Foo"}
	c3.LtmNode["Bar"] = &LtmNode{Name: "Bar"}
	c3.LtmRule["Foo"] = &LtmRule{Name: "Foo"}
	c3.LtmRule["Bar"] = &LtmRule{Name: "Bar"}

	err := c1.Merge(c2)
	if err != nil {
		t.Errorf("Error while merging F5 configuration: %s", err)
	}
	if !reflect.DeepEqual(c1, c3) {
		t.Errorf("Failed to merge F5 configuration")
	}
}

func TestParse(t *testing.T) {
	files := getFiles(testsF5config)
	repr.Println(files)
	cfg, err := ParseFile(testLog, files)
	if err != nil {
		t.Errorf("Cannot parse files: %s", err)
	}
	if false || true {
		repr.Println(cfg)
	}
	repr.Println(cfg.Info())
}
