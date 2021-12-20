package f5

import (
	"strings"

	"github.com/alecthomas/participle/v2/lexer"
)

// See https://clouddocs.f5.com/cli/tmsh-reference/latest/modules/ltm/ltm_rule.html
type LtmRule struct {
	OriginalConfig ParsedConfig
	Pos            lexer.Position
	Name           string `("ltm" "rule" @( F5Name | QF5Name ) | "rule" @( F5Name | QF5Name | Ident)) "{"?`
}

// newLtmRule parses data and creates a new LtmRule struct.
func newLtmRule(data ParsedConfig) (ret *LtmRule, err error) {
	ret = &LtmRule{}
	rule := strings.Split(data.Content, "\n")
	err = parseString("", rule[0], ret)
	ret.OriginalConfig = data
	return
}

func (o *LtmRule) Original() string {
	return o.OriginalConfig.Content
}
