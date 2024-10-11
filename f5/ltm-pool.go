// Copyright © 2023 Sébastien Gross <seb•ɑƬ•chezwam•ɖɵʈ•org>
//
// Created: 2021-12-19
// Last changed: 2024-10-11 21:18:28
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
	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
)

// See https://clouddocs.f5.com/cli/tmsh-reference/latest/modules/ltm/ltm_pool.html
type LtmPool struct {
	OriginalConfig         ParsedConfig
	Pos                    lexer.Position
	Name                   string           `("ltm" "pool" ( @F5Name | @QF5Name | @Ident ) | "pool" ( @F5Name | @QF5Name | @Ident)) "{"` // + version 10.x
	AllowNat               string           `(  "allow-nat" @( "yes" | "no" )`
	AllowSnat              string           ` | "allow-snat" @( "yes" | "no" )`
	AppService             string           ` | "app-service" @( "none" | QString | F5Name | QF5Name | Ident )`
	AutoscaleGroupID       string           ` | "autoscale-group-id" @( "none" | QString | Ident )`
	Description            string           ` | "description" @( QString | Ident )`
	GatewayFailsafeDevice  string           ` | "gateway-failsafe-device" @( QString | Ident )`
	IgnorePersistedWeight  string           ` | "ignore-persisted-weight" @( "yes" | "no" | "enabled" )`
	LoadBalancingMode      string           ` | "load-balancing-mode" @( F5lbMode  )` // TODO: limit algos
	Members                *[]LtmPoolMember ` | "members" ("{" @@* "}" | @@ ) `
	MinActiveMembers       int              ` | "min-active-members" @Ident`
	MinUpMembers           int              ` | "min-up-members" @Ident`
	MinUpMembersAction     string           ` | "min-up-members-action" @Ident`
	MinUpMembersChecking   string           ` | "min-up-members-checking" @Ident`
	Monitor                []string         ` | ("monitor" ( ( "min" Ident "of" "{" @( F5Name | QF5Name ) ("and"? @(F5Name | QF5Name) )* "}" ) | @( "none" | F5Name | QF5Name | (?!"all") Ident ) ("and"? @(F5Name | QF5Name | Ident) )* )) | ( "monitor" "all" @( Ident ) ("and" @( Ident ) )*  )` // all and indent for v10.x |  (?!"all") Ident -> do not match all ...
	Profiles               string           ` | "profiles" @( "none" | @F5Name | @QF5Name )`
	QueueOnConnectionLimit string           ` | "queue-on-connection-limit" @( "enabled" | "disabled" )`
	QueueDepthLimit        int              ` | "queue-depth-limit" @Ident`
	QueueTimeLimit         int              ` | "queue-time-limit" @Ident`
	ReselectTries          int              ` | "reselect-tries" @Ident`
	ServiceDownAction      string           ` | "service-down-action" @Ident`
	SlowRampTime           int              ` | ("slow-ramp-time" | "slow" "ramp" "time") @Ident  )* "}"`
}

type LtmPoolMember struct {
	Name            string        `@(F5Name | QF5Name | Ident) "{"`
	Address         string        ` ( "address" @Ident`
	AppService      string        ` | "app-service" @( "none" | QString | F5Name | QF5Name | Ident )`
	ConnectionLimit int           ` | "connection-limit" @Ident`
	Description     string        ` | "description" @( QString | Ident )`
	DynamicRatio    int           ` | "dynamic-ratio" @Ident`
	InheritProfile  *bool         ` | "inherit-profile" @("enabled" | "disabled")`
	Logging         *bool         ` | "logging" @("enabled" | "disabled")`
	Metadata        []*F5Metadata ` | "metadata" "{" @@  "}"`
	Monitor         []string      ` | "monitor" ( ( "min" Ident "of" "{" @( F5Name | QF5Name ) ("and"? @(F5Name | QF5Name) )* "}" ) | @( "none" | F5Name | QF5Name | Ident) ("and"? @(F5Name | QF5Name | Ident) )* )`
	PriorityGroup   string        ` | "priority-group" @( "none" | Ident )`
	RateLimit       int           ` | "rate-limit" @Ident`
	Ratio           int           ` | "ratio" @Ident`
	Session         string        ` | "session" @( "user-enabled" | "user-disabled" | "user" "disabled" | "monitor-enabled")`
	State           string        ` | "state" @( "user-up" | "user-down" | "down" | "up")`
	FQDN            *LtmPoolFQDN  ` | "fqdn" "{" @@ "}" )* "}"`
}

type LtmPoolFQDN struct {
	Autopopulate string `(  "autopopulate" @("enabled" | "disabled")`
	Name         string ` | "name" @Ident)*`
}

var ltmPoolParser = participle.MustBuild[LtmPool](
	participle.Lexer(f5Lexer),
	participle.Unquote("QF5Name"),
	participle.Unquote("QString"),
)

// newLtmPool parses data and creates a new LtmPool struct.
func newLtmPool(data ParsedConfig) (ret *LtmPool, err error) {
	ret = &LtmPool{}
	ret, err = ltmPoolParser.ParseString("", data.Content)
	ret.OriginalConfig = data
	return
}

func (o *LtmPool) Original() string {
	return o.OriginalConfig.Content
}

func (o *LtmPool) GetName() string {
	return o.Name
}
