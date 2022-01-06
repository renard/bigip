package f5

import (
	"embed"
	"testing"

	"github.com/alecthomas/repr"
)

var (
	//go:embed testdata/ltm-monitor-*.conf
	testsLtmMonitor embed.FS
)

func TestLtmMonitor(t *testing.T) {
	files := getFiles(testsLtmMonitor)
	for _, file := range files {
		data, _ := testsLtmMonitor.ReadFile(file)
		obj, err := newLtmMonitor(ParsedConfig{Content: string(data)})
		if err != nil {
			t.Errorf("%s Cannot parse monitor snippet: %s", file, err)
		}
		if false {
			repr.Println(obj)
		}
	}
}
