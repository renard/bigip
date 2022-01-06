package f5

import (
	"embed"
	"testing"

	"github.com/alecthomas/repr"
)

var (
	//go:embed testdata/ltm-pool-*.conf
	testsLtmPool embed.FS
)

func TestLtmPool(t *testing.T) {
	files := getFiles(testsLtmPool)
	for _, file := range files {
		data, _ := testsLtmPool.ReadFile(file)
		obj, err := newLtmPool(ParsedConfig{Content: string(data)})
		if err != nil {
			t.Errorf("%s Cannot parse virtual snippet: %s", file, err)
		}
		if false {
			repr.Println(obj)
		}
	}
}
