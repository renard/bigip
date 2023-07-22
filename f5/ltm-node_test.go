// Copyright © 2023 Sébastien Gross <seb•ɑƬ•chezwam•ɖɵʈ•org>
//
// Created: 2021-12-19
// Last changed: 2023-07-22 02:57:57
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
	"testing"

	"github.com/alecthomas/repr"
)

var (
	//go:embed testdata/ltm-node-*.conf
	testsLtmNode embed.FS
)

func TestLtmNode(t *testing.T) {
	files := getFiles(testsLtmNode)
	for _, file := range files {
		data, _ := testsLtmNode.ReadFile(file)
		obj, err := newLtmNode(ParsedConfig{Content: string(data)})
		if err != nil {
			t.Errorf("%s Cannot parse virtual snippet: %s", file, err)
		}
		if false {
			repr.Println(obj)
		}
	}
}
