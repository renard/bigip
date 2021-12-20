package f5

import (
	"embed"
	"strings"
	"testing"

	"github.com/alecthomas/repr"
)

var (
	//go:embed testdata/ltm-virtual-*.cfg
	testsLtmVirtual embed.FS
)

func TestLtmVirtual(t *testing.T) {
	files := getFiles(testsLtmVirtual)
	for _, file := range files {
		data, _ := testsLtmVirtual.ReadFile(file)
		obj, err := newLtmVirtual(ParsedConfig{Content: strings.Split(string(data), "\n")})
		if err != nil {
			t.Errorf("%s Cannot parse virtual snippet: %s", file, err)
		}
		if false {
			repr.Println(obj)
		}
	}
}
