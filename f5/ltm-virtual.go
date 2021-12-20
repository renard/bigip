package f5

import (
	"github.com/alecthomas/participle/v2/lexer"
)

// See https://clouddocs.f5.com/cli/tmsh-reference/latest/modules/ltm/ltm_virtual.html
type LtmVirtual struct {
	original                 string
	Pos                      lexer.Position
	Name                     string                                `"ltm" "virtual" ( @F5Name | @QF5Name ) "{"`
	Description              string                                `( "description" @( QString | Ident )`
	AddressStatus            string                                ` | "address-status" @( "yes" | "no" )`
	AppService               string                                ` | "app-service" @( "none" | QString | Ident )`
	AutoDiscovery            string                                ` | "auto-discovery" @("enabled" | "disabled")`
	AutoLasthop              string                                ` | "auto-lasthop" @( "default" | "enabled" | "disabled")`
	BwcPolicy                string                                ` | "bwc-policy" @( F5Name | QF5Name )` // OLD Devices?
	ConnectionLimit          int                                   ` | "connection-limit" @( Ident )`      // OLD Devices?
	ClonePools               []*LtmVirtualClonePool                ` | "clone-pools" "{" @@+ "}"`
	Destination              string                                ` | "destination" @( F5Name | QF5Name )`
	Enabled                  string                                ` | @( "enabled" | "disabled" )`
	FallbackPersistence      string                                ` | "fallback-persistence" @( F5Name | QF5Name )`
	IpProtocol               string                                ` | "ip-protocol" @( "any" | "udp" | "tcp" )`
	Mask                     string                                ` | "mask" @Ident`
	Pool                     string                                ` | "pool" @( F5Name | QF5Name )`
	Persist                  []*LtmVirtualPersist                  ` | "persist" "{" @@+ "}"`
	Profiles                 []*LtmVirtualProfile                  ` | "profiles" "{" @@+ "}"`
	Policies                 []*LtmVirtualPolicy                   ` | "policies" "{" @@+ "}"`
	Rules                    []string                              ` | "rules" "{" @( QF5Name | F5Name)+ "}"`
	SecurityLogProfiles      []string                              ` | "security-log-profiles" "{" @(Ident | QF5Name | F5Name | QString)* "}"` // OLD Devices?
	Source                   string                                ` | "source" @Ident`
	SourcePort               string                                ` | "source-port" @Ident`
	SourceAddressTranslation []*LtmVirtualSourceAddressTranslation ` | "source-address-translation" "{" @@ "}"`
	TranslateAddress         string                                ` | "translate-address" @( "enabled" | "disabled")`
	TranslatePort            string                                ` | "translate-port" @( "enabled" | "disabled")`
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
	Name    string `@( F5Name | QF5Name ) "{"`
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
	Name    string `@( F5Name | QF5Name ) "{"`
	Context string `("context" @( "all" | "clientside" | "serverside") )? "}"`
}

type LtmVirtualSourceAddressTranslation struct {
	Pool string `(  "pool" @( "none" | F5Name | QF5Name )`
	Type string ` | "type" @( "automap" | "lsn" | "snat" | "none" ))+`
}

// newLtmVirtual parses data and creates a new LtmVirtual struct.
func newLtmVirtual(data ParsedConfig) (ret *LtmVirtual, err error) {
	ret = &LtmVirtual{}
	err = parseString("", data.Content, ret)
	ret.original = data.Content
	return
}

func (o *LtmVirtual) Original() string {
	return o.original
}
