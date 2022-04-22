package f5

import (
	"embed"
	"testing"

	"github.com/alecthomas/repr"
)

var (
	//go:embed testdata/ltm-policy-*.conf
	testsLtmPolicy embed.FS
)

func TestLtmPolicy(t *testing.T) {
	files := getFiles(testsLtmPolicy)
	for _, file := range files {
		data, _ := testsLtmPolicy.ReadFile(file)
		obj, err := newLtmPolicy(ParsedConfig{Content: string(data)})
		if err != nil {
			t.Errorf("%s Cannot parse policy snippet: %s", file, err)
		}
		if false {
			repr.Println(obj)
		}
	}
}
