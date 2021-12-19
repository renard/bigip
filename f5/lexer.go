package f5

import (
	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
)

var (
	f5Lexer = lexer.MustSimple(
		[]lexer.Rule{
			{"comment", `#[^\n]*`, nil},
			{"whitespace", `\s+`, nil},
			{"eol", `[\n\r]+`, nil},
			{"Punct", `[{}]`, nil},
			// Name definition
			{"QF5Name", `"/[A-Za-z0-9/_. -]*[A-Za-z0-9/_%:]+"`, nil},
			{"F5Name", `/[A-Za-z0-9/_.-]*[A-Za-z0-9/_%:]+`, nil},
			{"F5Time", `\d{4}-\d{2}-\d{2}:\d{2}:\d{2}:\d{2}`, nil},
			// {"F5IPCIDR", `\b[0-9a-fA-F\.:]+/[0-9]+\b`, nil},
			// {"F5IP", `\b[0-9a-fA-F\.:]+\b`, nil},

			{"F5lbMode", `dynamic-ratio-member|dynamic-ratio-node|fastest-app-response|fastest-node|least-connections-members?|least-connections-node|least-sessions|observed-member|observed-node|predictive-member|predictive-node|ratio-least-connections-member|ratio-least-connections-node|ratio-member|ratio-node|ratio-session|round-robin|weighted-least-connections-member|weighted-least-connections-node`, nil},
			// Quoted string
			{"QString", `".+"`, nil},

			{"Ident", `[A-Za-z0-9._-][A-Za-z0-9._/:^%-]*`, nil},
		},
		lexer.MatchLongest(),
	)
)

func parseString(name, data string, obj interface{}) (err error) {
	parser := participle.MustBuild(
		obj,
		participle.Lexer(f5Lexer),
		participle.Unquote("QF5Name"),
		participle.Unquote("QString"),
	)
	err = parser.ParseString(name, data, obj)
	return
}
