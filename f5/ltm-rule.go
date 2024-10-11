// Copyright © 2023 Sébastien Gross <seb•ɑƬ•chezwam•ɖɵʈ•org>
//
// Created: 2021-12-20
// Last changed: 2024-10-11 21:20:17
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

// See https://clouddocs.f5.com/cli/tmsh-reference/latest/modules/ltm/ltm_rule.html
type LtmRule struct {
	OriginalConfig ParsedConfig
	Pos            lexer.Position
	Name           string `("ltm" "rule" @( F5Name | QF5Name ) | "rule" @( F5Name | QF5Name | Ident)) "{"?`
}

var ltmRuleParser = participle.MustBuild[LtmRule](
	participle.Lexer(f5Lexer),
	participle.Unquote("QF5Name"),
	participle.Unquote("QString"),
)

// newLtmRule parses data and creates a new LtmRule struct.
func newLtmRule(data ParsedConfig) (ret *LtmRule, err error) {
	ret = &LtmRule{}
	rule := strings.Split(data.Content, "\n")
	ret, err = ltmRuleParser.ParseString("", rule[0])
	ret.OriginalConfig = data
	return
}

func (o *LtmRule) Original() string {
	return o.OriginalConfig.Content
}

func (o *LtmRule) GetName() string {
	return o.Name
}
