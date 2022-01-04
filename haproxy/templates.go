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

func GenerateTemplates(config *Config, cfg f5.F5Config) (err error) {
	tmpls, err := loadTemplates(config)
	if err != nil {
		return
	}

	t := tmpls.Lookup("export")
	err = t.Execute(os.Stdout, struct {
		Cfg f5.F5Config
	}{
		Cfg: cfg,
	})
	return
}
