package haproxy

import (
	"bytes"
	"embed"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"bigip/f5"

	"github.com/Masterminds/sprig"
	"github.com/alecthomas/repr"
)

var (
	//go:embed templates/*
	tpl embed.FS
)

type Config struct {
	TemplateDir []string
	Export      []string
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

func loadTemplates(config *Config) (tmpls *template.Template, err error) {
	tmpls = template.New("")
	funcs := template.FuncMap{
		"comment":  comment,
		"scomment": spacedComment,
		"ipport":   ipport,
		// https://forum.golangbridge.org/t/template-check-if-block-is-defined/6928/2
		"hasTemplate": func(name string) bool {
			return tmpls.Lookup(name) != nil
		},
		"templateIfExists": func(name string, pipeline interface{}) (string, error) {
			t := tmpls.Lookup(name)
			if t == nil {
				return "", nil
			}

			buf := &bytes.Buffer{}
			err := t.Execute(buf, pipeline)
			if err != nil {
				return "", err
			}

			return buf.String(), nil
		},
	}
	tmpls = tmpls.Funcs(sprig.TxtFuncMap())
	tmpls = tmpls.Funcs(funcs)
	tmpls, err = tmpls.ParseFS(tpl, "templates/*.cfg")
	if err != nil {
		return
	}

	for _, td := range config.TemplateDir {
		err = addExtraTemplates(tmpls, td)
		if err != nil {
			return
		}
	}
	return
}

func Render(config *Config, cfg f5.F5Config) (err error) {
	tmpls, err := loadTemplates(config)
	if err != nil {
		return
	}

	if false {
		repr.Println(tmpls)
	}

	t := tmpls.Lookup("main")
	err = t.Execute(os.Stdout, struct {
		Cfg f5.F5Config
	}{
		Cfg: cfg,
	})
	return
}

func GenerateTemplates(config *Config, f5config f5.F5Config) (err error) {
	tmpls, err := loadTemplates(config)
	if err != nil {
		return
	}

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
			}
		}
	}
	t := tmpls.Lookup("export")
	err = t.Execute(os.Stdout, struct {
		F5config f5.F5Config
		Config   Config
	}{
		F5config: f5c,
		Config:   *config,
	})
	return
}
