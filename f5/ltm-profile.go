// Copyright © 2023 Sébastien Gross <seb•ɑƬ•chezwam•ɖɵʈ•org>
//
// Created: 2021-12-20
// Last changed: 2024-10-11 21:19:18
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
	"strings"

	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
)

// See https://clouddocs.f5.com/cli/tmsh-reference/latest/modules/ltm/ltm_profile.html
type LtmProfile struct {
	OriginalConfig ParsedConfig
	Pos            lexer.Position
	Type           string `"ltm" "profile" @Ident`
	Name           string `( @F5Name | @QF5Name ) "{"?`
}

var ltmProfileParser = participle.MustBuild[LtmProfile](
	participle.Lexer(f5Lexer),
	participle.Unquote("QF5Name"),
	participle.Unquote("QString"),
)

// newLtmProfile parses data and creates a new LtmProfile struct.
func newLtmProfile(data ParsedConfig) (ret *LtmProfile, err error) {
	ret = &LtmProfile{}
	profile := strings.Split(data.Content, "\n")
	ret, err = ltmProfileParser.ParseString("", profile[0])
	ret.OriginalConfig = data
	return
}

func (o *LtmProfile) Original() string {
	return o.OriginalConfig.Content
}

func (o *LtmProfile) GetName() string {
	return o.Name
}
