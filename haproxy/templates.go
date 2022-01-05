package haproxy

import (
	"embed"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"bigip/f5"

	"github.com/alecthomas/repr"
)

var (
	//go:embed templates/*
	tpl embed.FS
)

type Config struct {
	TemplateDir []string
	Export      []string
	OutputDir   string
}

func addExtraTemplates(t *template.Template, dir string) (err error) {
	err = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if strings.Contains(path, ".cfg") {
			_, err = t.ParseFiles(path)
			if err != nil {
				repr.Println(err)
			}
		}
		return err
	})
	return err
}

func Render(config *Config, f5config f5.F5Config) (err error) {
	tmpls, err := loadTemplates(config)
	if err != nil {
		return
	}

	if false {
		repr.Println(tmpls)
	}

	t := tmpls.Lookup("main")
	err = t.Execute(os.Stdout, struct {
		F5config f5.F5Config
		Config   Config
	}{
		F5config: f5config,
		Config:   *config,
	})
	return
}

func GenerateTemplates(config *Config, f5config f5.F5Config) (err error) {
	tmpls, err := loadTemplates(config)
	if err != nil {
		return
	}

	err = os.MkdirAll(config.OutputDir, os.ModePerm)
	if err != nil {
		return
	}

	t := tmpls.Lookup("export")
	for _, tp := range []string{"rule", "profile", "node", "monitor"} {
		f5c := f5.NewF5Config()
		switch tp {
		case "rule":
			f5c.LtmRule = f5config.LtmRule
		case "profile":
			f5c.LtmProfile = f5config.LtmProfile
		case "node":
			f5c.LtmNode = f5config.LtmNode
		case "monitor":
			f5c.LtmMonitor = f5config.LtmMonitor
		}

		fh, e := os.Create(fmt.Sprintf("%s/%s.cfg", config.OutputDir, tp))
		if e != nil {
			return e
		}

		err = t.Execute(fh, struct {
			F5config f5.F5Config
			Config   Config
		}{
			F5config: f5c,
			Config:   *config,
		})
		if err != nil {
			return
		}
		fh.Close()
	}

	if false {
		f5c := f5.NewF5Config()
		if len(config.Export) == 0 {
			f5c = f5config
		} else {
			for _, f := range config.Export {
				switch f {
				case "rule":
					f5c.LtmRule = f5config.LtmRule
				case "profile":
					f5c.LtmProfile = f5config.LtmProfile
				case "node":
					f5c.LtmNode = f5config.LtmNode
				case "monitor":
					f5c.LtmMonitor = f5config.LtmMonitor
				}
			}
		}
		//t := tmpls.Lookup("export")
		err = t.Execute(os.Stdout, struct {
			F5config f5.F5Config
			Config   Config
		}{
			F5config: f5c,
			Config:   *config,
		})
	}
	return
}
