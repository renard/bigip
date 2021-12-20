package f5

import (
	"embed"
	"testing"

	"github.com/alecthomas/repr"
)

var (
	//go:embed testdata/ltm-profile-*.cfg
	testsLtmProfile embed.FS
)

func TestLtmProfile(t *testing.T) {
	files := getFiles(testsLtmProfile)
	for _, file := range files {
		data, _ := testsLtmProfile.ReadFile(file)
		obj, err := newLtmProfile(ParsedConfig{Content: string(data)})
		if err != nil {
			t.Errorf("%s Cannot parse profile snippet: %s", file, err)
		}
		if false || true {
			repr.Println(obj)
		}
	}
}
