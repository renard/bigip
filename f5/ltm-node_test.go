package f5

import (
	"embed"
	"testing"

	"github.com/alecthomas/repr"
)

var (
	//go:embed testdata/ltm-node-*.conf
	testsLtmNode embed.FS
)

func TestLtmNode(t *testing.T) {
	files := getFiles(testsLtmNode)
	for _, file := range files {
		data, _ := testsLtmNode.ReadFile(file)
		obj, err := newLtmNode(ParsedConfig{Content: string(data)})
		if err != nil {
			t.Errorf("%s Cannot parse virtual snippet: %s", file, err)
		}
		if false {
			repr.Println(obj)
		}
	}
}
