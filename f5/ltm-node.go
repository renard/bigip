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
