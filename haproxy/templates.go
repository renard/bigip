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
	Files           []string
	TemplateDir     []string
	Export          []string
	OutputDir       string
	ExpandTemplates bool
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

	out := os.Stdout
	if config.OutputDir != "" && config.OutputDir != "-" {
		out, err = os.Create(fmt.Sprintf("%s/main.cfg", config.OutputDir))
		defer out.Close()
	}

	t := tmpls.Lookup("main")
	err = t.Execute(out, struct {
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
	for _, tp := range []string{"virtual", "pool", "rule", "profile", "node", "monitor",
		"persistence"} {
		f5c := f5.NewF5Config()
		switch tp {
		case "virtual":
			f5c.LtmVirtual = f5config.LtmVirtual
		case "pool":
			f5c.LtmPool = f5config.LtmPool
			f5c.LtmNode = f5config.LtmNode
		case "rule":
			f5c.LtmRule = f5config.LtmRule
		case "profile":
			f5c.LtmProfile = f5config.LtmProfile
		case "node":
			f5c.LtmNode = f5config.LtmNode
		case "monitor":
			f5c.LtmMonitor = f5config.LtmMonitor
		case "persistence":
			f5c.LtmPersistence = f5config.LtmPersistence
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
	return
}
