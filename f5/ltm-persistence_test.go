package f5

import (
	"embed"
	"testing"

	"github.com/alecthomas/repr"
)

var (
	//go:embed testdata/ltm-persistence-*.conf
	testsLtmPersistence embed.FS
)

func TestLtmPersistence(t *testing.T) {
	files := getFiles(testsLtmPersistence)
	for _, file := range files {
		data, _ := testsLtmPersistence.ReadFile(file)
		obj, err := newLtmPersistence(ParsedConfig{Content: string(data)})
		if err != nil {
			t.Errorf("%s Cannot parse profile snippet: %s", file, err)
		}
		if false {
			repr.Println(obj)
		}
	}
}
