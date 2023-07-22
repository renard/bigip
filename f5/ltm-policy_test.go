// Copyright © 2023 Sébastien Gross <seb•ɑƬ•chezwam•ɖɵʈ•org>
//
// Created: 2022-04-22
// Last changed: 2023-07-22 02:58:16
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
	//go:embed testdata/ltm-policy-*.conf
	testsLtmPolicy embed.FS
)

func TestLtmPolicy(t *testing.T) {
	files := getFiles(testsLtmPolicy)
	for _, file := range files {
		data, _ := testsLtmPolicy.ReadFile(file)
		obj, err := newLtmPolicy(ParsedConfig{Content: string(data)})
		if err != nil {
			t.Errorf("%s Cannot parse policy snippet: %s", file, err)
		}
		if false {
			repr.Println(obj)
		}
	}
}
