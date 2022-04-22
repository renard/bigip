package f5

import (
	"strings"

	"github.com/alecthomas/participle/v2/lexer"
)

// See https://clouddocs.f5.com/cli/tmsh-reference/latest/modules/ltm/ltm_policy.html
type LtmPolicy struct {
	OriginalConfig ParsedConfig
	Pos            lexer.Position
	Name           string `"ltm" "policy" ( @F5Name | @QF5Name ) "{"?`
}

// newLtmPolicy parses data and creates a new LtmPolicy struct.
func newLtmPolicy(data ParsedConfig) (ret *LtmPolicy, err error) {
	ret = &LtmPolicy{}
	profile := strings.Split(data.Content, "\n")
	err = parseString("", profile[0], ret)
	ret.OriginalConfig = data
	return
}

func (o *LtmPolicy) Original() string {
	return o.OriginalConfig.Content
}

func (o *LtmPolicy) GetName() string {
	return o.Name
}
