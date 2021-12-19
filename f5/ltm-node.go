package f5

import (
	"github.com/alecthomas/participle/v2/lexer"
)

// See https://clouddocs.f5.com/cli/tmsh-reference/latest/modules/ltm/ltm_node.html
type LtmNode struct {
	original        string
	Pos             lexer.Position
	Name            string       `"ltm" "node" ( @F5Name | @QF5Name ) "{"`
	Description     string       `( "description" @( QString | Ident )`
	Address         string       ` | "address" @Ident`
	AppService      string       ` | "app-service" @( "none" | QString | Ident )`
	ConnectionLimit int          ` | "connection-limit" @Ident`
	Status          string       ` | @( "up" | "down" )`
	DynamicRatio    int          ` | "dynamic-ratio" @Ident`
	Fqdn            *LtmNodeFQDN ` | "fqdn" "{" @@ "}"`
	Logging         *bool        ` | "logging" @("enabled" | "disabled")`
	Monitor         string       ` | "monitor" @( "none" | @F5Name | @QF5Name )`
	RateLimit       int          ` | "rate-limit" @Ident`
	Ratio           int          ` | "ratio" @Ident`
	Session         string       ` | @( "user-enabled" | "user-disabled")`
	State           string       ` | @( "user-up" | "user-down" ) )* "}"`
}

type LtmNodeFQDN struct {
	AddressFamily string `( "address-family" @Ident`
	Autopopulate  *bool  ` | "autopopulate" @("enabled" | "disabled")`
	Name          string ` | "name" @Ident`
	Interval      int    ` | "interval" @Ident`
	DownInterval  int    ` | "down-interval" @Ident )*`
}

// newLtmNode parses data and creates a new LtmNode struct.
func newLtmNode(data string) (ret *LtmNode, err error) {
	ret = &LtmNode{}
	err = parseString("", data, ret)
	ret.original = data
	return
}

func (o *LtmNode) Original() string {
	return o.original
}
