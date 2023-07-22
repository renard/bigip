// Copyright © 2023 Sébastien Gross <seb•ɑƬ•chezwam•ɖɵʈ•org>
//
// Created: 2021-12-19
// Last changed: 2023-07-22 02:57:52
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
	"github.com/alecthomas/participle/v2/lexer"
)

// See https://clouddocs.f5.com/cli/tmsh-reference/latest/modules/ltm/ltm_node.html
type LtmNode struct {
	OriginalConfig  ParsedConfig
	Pos             lexer.Position
	Name            string       `("ltm" "node" ( @F5Name | @QF5Name ) | "node" ( @F5Name | @QF5Name | @Ident) )"{"` // + version 10.x
	Description     string       `( "description" @( QString | Ident )`
	Address         string       ` | "address" @Ident`
	AppService      string       ` | "app-service" @( "none" | QString | Ident )`
	ConnectionLimit int          ` | "connection-limit" @Ident`
	Status          string       ` | @( "up" | "down" )`
	DynamicRatio    int          ` | "dynamic-ratio" @Ident`
	Fqdn            *LtmNodeFQDN ` | "fqdn" "{" @@ "}"`
	Logging         *bool        ` | "logging" @("enabled" | "disabled")`
	Monitor         []string     ` | "monitor" @( ("none" | F5Name | QF5Name | Ident ) ("and"? ( F5Name | QF5Name | Ident ))*)` // @Ident for version 10.x
	RateLimit       int          ` | "rate-limit" @Ident`
	Ratio           int          ` | "ratio" @Ident`
	Session         string       ` | "session" @( "user-enabled" | "user-disabled")`
	Screen          string       ` | "screen" @Ident` // version 10.x
	State           string       ` | "state" @( "user-up" | "user-down" ) )* "}"`
}

type LtmNodeFQDN struct {
	AddressFamily string `( "address-family" @Ident`
	Autopopulate  *bool  ` | "autopopulate" @("enabled" | "disabled")`
	Name          string ` | "name" @Ident`
	Interval      int    ` | "interval" @Ident`
	DownInterval  int    ` | "down-interval" @Ident )*`
}

// newLtmNode parses data and creates a new LtmNode struct.
func newLtmNode(data ParsedConfig) (ret *LtmNode, err error) {
	ret = &LtmNode{}
	err = parseString("", data.Content, ret)
	ret.OriginalConfig = data
	return
}

func (o *LtmNode) Original() string {
	return o.OriginalConfig.Content
}

func (o *LtmNode) GetName() string {
	return o.Name
}
