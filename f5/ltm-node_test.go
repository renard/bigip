package f5

import (
	"embed"
	"strings"
	"testing"

	"github.com/alecthomas/repr"
)

var (
	//go:embed testdata/ltm-node-*.cfg
	testsLtmNode embed.FS
)

func TestLtmNode(t *testing.T) {
	files := getFiles(testsLtmNode)
	for _, file := range files {
		data, _ := testsLtmNode.ReadFile(file)
		obj, err := newLtmNode(ParsedConfig{Content: strings.Split(string(data), "\n")})
		if err != nil {
			t.Errorf("%s Cannot parse virtual snippet: %s", file, err)
		}
		if false {
			repr.Println(obj)
		}
	}
}
