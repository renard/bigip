package f5

import (
	"strings"

	"github.com/alecthomas/participle/v2/lexer"
)

// See https://clouddocs.f5.com/cli/tmsh-reference/latest/modules/ltm/ltm_persistence_cookie.html
type LtmPersistence struct {
	OriginalConfig ParsedConfig
	Pos            lexer.Position
	Type           string `"ltm" "persistence" @Ident`
	Name           string `( @F5Name | @QF5Name ) "{"?`
}

// newLtmPersistence parses data and creates a new LtmPersistence struct.
func newLtmPersistence(data ParsedConfig) (ret *LtmPersistence, err error) {
	ret = &LtmPersistence{}
	profile := strings.Split(data.Content, "\n")
	err = parseString("", profile[0], ret)
	ret.OriginalConfig = data
	return
}

func (o *LtmPersistence) Original() string {
	return o.OriginalConfig.Content
}
