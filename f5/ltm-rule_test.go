package f5

import (
	"embed"
	"testing"

	"github.com/alecthomas/repr"
)

var (
	//go:embed testdata/ltm-rule-*.conf
	testsLtmRule embed.FS
)

func TestLtmRule(t *testing.T) {
	files := getFiles(testsLtmRule)
	for _, file := range files {
		data, _ := testsLtmRule.ReadFile(file)
		obj, err := newLtmRule(ParsedConfig{Content: string(data)})
		if err != nil {
			t.Errorf("%s Cannot parse rule snippet: %s", file, err)
		}
		if false {
			repr.Println(obj)
		}
	}
}
