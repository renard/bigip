// Copyright © 2023 Sébastien Gross <seb•ɑƬ•chezwam•ɖɵʈ•org>
//
// Created: 2021-12-19
// Last changed: 2024-10-11 21:23:36
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
	// "github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
)

var (
	f5Lexer = lexer.MustSimple(
		[]lexer.SimpleRule{
			{"comment", `#[^\n]*`},
			{"whitespace", `\s+`},
			{"eol", `[\n\r]+`},
			{"Punct", `[{}]`},
			// Name definition
			{"QF5Name", `"/[A-Za-z0-9/_. -]*[A-Za-z0-9/_%:]+"`},
			{"F5Name", `/[A-Za-z0-9/_.-]*[A-Za-z0-9/_%:]+`},
			{"F5Time", `\d{4}-\d{2}-\d{2}:\d{2}:\d{2}:\d{2}`},
			// {"F5IPCIDR", `\b[0-9a-fA-F\.:]+/[0-9]+\b`},
			// {"F5IP", `\b[0-9a-fA-F\.:]+\b`},

			{"F5lbMode", `dynamic-ratio-member|dynamic-ratio-node|fastest-app-response|fastest-node|least-connections-members?|least-connections-node|least-sessions|observed-member|observed-node|predictive-member|predictive-node|ratio-least-connections-member|ratio-least-connections-node|ratio-member|ratio-node|ratio-session|round-robin|weighted-least-connections-member|weighted-least-connections-node`},
			// Quoted string
			{"QString", `".+"`},

			{"Ident", `[A-Za-z0-9._-][A-Za-z0-9._/:^%-]*`},
		},
		// lexer.MatchLongest(),
	)
)
