package f5

import (
	"strings"

	"github.com/alecthomas/participle/v2/lexer"
)

// See https://clouddocs.f5.com/cli/tmsh-reference/latest/modules/ltm/ltm_profile.html
type LtmProfile struct {
	OriginalConfig ParsedConfig
	Pos            lexer.Position
	Type           string `"ltm" "profile" @Ident`
	Name           string `( @F5Name | @QF5Name ) "{"?`
}

// newLtmProfile parses data and creates a new LtmProfile struct.
func newLtmProfile(data ParsedConfig) (ret *LtmProfile, err error) {
	ret = &LtmProfile{}
	profile := strings.Split(data.Content, "\n")
	err = parseString("", profile[0], ret)
	ret.OriginalConfig = data
	return
}

func (o *LtmProfile) Original() string {
	return o.OriginalConfig.Content
}
