// Copyright © 2023 Sébastien Gross <seb•ɑƬ•chezwam•ɖɵʈ•org>
//
// Created: 2021-12-19
// Last changed: 2024-10-11 01:16:53
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

// See https://clouddocs.f5.com/cli/tmsh-reference/latest/modules/ltm/ltm_virtual.html
type LtmVirtual struct {
	OriginalConfig             ParsedConfig
	Pos                        lexer.Position
	Name                       string                                `("ltm" "virtual" @(F5Name | QF5Name | Ident ) | "virtual" "address"? @(F5Name | QF5Name | Ident)) "{"`
	Description                string                                `( "description" @( QString | Ident )`
	AddressStatus              string                                ` | "address-status" @( "yes" | "no" )`
	AppService                 string                                ` | "app-service" @( "none" | QString | F5Name | QF5Name | Ident )`
	AutoDiscovery              string                                ` | "auto-discovery" @("enabled" | "disabled")`
	AutoLasthop                string                                ` | "auto-lasthop" @( "default" | "enabled" | "disabled")`
	BwcPolicy                  string                                ` | "bwc-policy" @( F5Name | QF5Name )` // OLD Devices?
	ConnectionLimit            int                                   ` | "connection-limit" @( Ident )`      // OLD Devices?
	ClonePools                 []*LtmVirtualClonePool                ` | "clone-pools" "{" @@+ "}"`
	Destination                string                                ` | "destination" @( F5Name | QF5Name | Ident )`
	Enabled                    string                                ` | @( "enabled" | "disabled" | "disable" )`
	FallbackPersistence        string                                ` | "fallback-persistence" @( F5Name | QF5Name )`
	IpProtocol                 string                                ` | ("ip-protocol" | "ip" "protocol") @( "any" | "udp" | "tcp" )`
	IpForward                  *bool                                 ` | @"ip-forward"`
	Mask                       string                                ` | "mask" @Ident`
	Pool                       string                                ` | "pool" @( F5Name | QF5Name | Ident )`
	Persist                    []*LtmVirtualPersist                  ` | "persist" "{" @@+ "}"`
	PersistOLD                 string                                ` | "persist" @Ident`
	Profiles                   []*LtmVirtualProfile                  ` | "profiles" "{" @@+ "}"`
	Policies                   []*LtmVirtualPolicy                   ` | "policies" "{" @@+ "}"`
	Rules                      []string                              ` | "rules" ( "{" @( QF5Name | F5Name | Ident)+ "}" | @( QF5Name | F5Name | Ident) )`
	SecurityLogProfiles        []string                              ` | "security-log-profiles" "{" @(Ident | QF5Name | F5Name | QString)* "}"` // OLD Devices?
	ServerSSLUseSNI            string                                ` | "serverssl-use-sni"  @( Ident )`                                        // OLD Devices?
	ServiceDownImmediateAction string                                ` | "service-down-immediate-action"  @( Ident )`                            // OLD Devices?
	Source                     string                                ` | "source" @Ident`
	SourcePort                 string                                ` | "source-port" @Ident`
	SourceAddressTranslation   []*LtmVirtualSourceAddressTranslation ` | "source-address-translation" "{" @@ "}"`
	TrafficMatchingCriteria    string                                ` | "traffic-matching-criteria"  @(Ident | QF5Name | F5Name | QString)`
	Translate                  string                                ` | "translate" "service" @( "enable" | "disable")` // OLD device
	TranslateAddress           string                                ` | "translate-address" @( "enabled" | "disabled")`
	TranslatePort              string                                ` | "translate-port" @( "enabled" | "disabled")`
	Vlans                      []string                              ` | "vlans" "{" @(Ident | QF5Name | F5Name | QString)* "}"` // OLD Devices?
	VlansEnabled               *bool                                 ` | @"vlans-enabled"`                                       // OLD Devices?
	VsIndex                    int                                   ` | "vs-index" @Ident`                                      // OLD Devices?
	Metadata                   []*F5Metadata                         ` | "metadata" "{" @@* "}"`
	// Address         string          ` | "address" @Ident`
	// ConnectionLimit int             ` | "connection-limit" @Ident`
	// Status          string          ` | @( "up" | "down" )`
	// DynamicRatio    int             ` | "dynamic-ratio" @Ident`
	// Fqdn            *LtmVirtualFQDN ` | "fqdn" "{" @@ "}"`
	// Logging         *bool           ` | "logging" @("enabled" | "disabled")`
	// Monitor         string          ` | "monitor" @( "none" | @F5Name | @QF5Name )`
	// RateLimit       int             ` | "rate-limit" @Ident`
	// Ratio           int             ` | "ratio" @Ident`
	// Session         string          ` | @( "user-enabled" | "user-disabled")`
	CreationTime     string ` | "creation-time" @F5Time`
	LastModifiedTime string ` | "last-modified-time" @F5Time )* "}"`
	// State string ` | @( "user-up" | "user-down" ) )* "}"`
}

type LtmVirtualPersist struct {
	Name    string `@( F5Name | QF5Name | Ident ) "{"`
	Default string `("default" @( "yes" | "no" ) )? "}"`
}

type LtmVirtualPolicy struct {
	Name    string `@( F5Name | QF5Name ) "{"`
	Context string `("context" @( "all" | "clientside" | "serverside") )? "}"`
}
type LtmVirtualClonePool struct {
	Name    string `@( F5Name | QF5Name ) "{"`
	Context string `("context" @( "all" | "clientside" | "serverside") )? "}"`
}

type LtmVirtualProfile struct {
	Name    string `@( F5Name | QF5Name | Ident ) "{"`
	Context string `("context" @( "all" | "clientside" | "serverside") )? "}"`
}

type LtmVirtualSourceAddressTranslation struct {
	Pool string `(  "pool" @( "none" | F5Name | QF5Name | Ident )`
	Type string ` | "type" @( "automap" | "lsn" | "snat" | "none" ))+`
}

// newLtmVirtual parses data and creates a new LtmVirtual struct.
func newLtmVirtual(data ParsedConfig) (ret *LtmVirtual, err error) {
	ret = &LtmVirtual{}
	err = parseString("", data.Content, ret)
	ret.OriginalConfig = data
	return
}

func (o *LtmVirtual) Original() string {
	return o.OriginalConfig.Content
}

func (o *LtmVirtual) GetName() string {
	return o.Name
}
