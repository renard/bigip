package haproxy

import (
	"bytes"
	"embed"
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

func comment(str string, indent int) string {
	idstr := strings.Repeat(" ", indent)
	lines := strings.Split(str, "\n")
	for i, line := range lines {
		lines[i] = idstr + "# " + line
	}
	return strings.Join(lines, "\n")
}

func ipport(str string) string {
	strs := strings.Split(str, "/")
	return strs[len(strs)-1]
}

// func hasTemplate(name string) bool {
// 	return tmpl.Lookup(name) != nil
// }

// func templateIfExists(name string, pipeline interface{}) (string, error) {
// 	t := tmpl.Lookup(name)
// 	if t == nil {
// 		return "", nil
// 	}

// 	buf := &bytes.Buffer{}
// 	err := t.Execute(buf, pipeline)
// 	if err != nil {
// 		return "", err
// 	}

// 	return buf.String(), nil
// }

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

// https://forum.golangbridge.org/t/template-check-if-block-is-defined/6928/2
func Render(cfg f5.F5Config) (err error) {
	tmpls := template.New("")
	funcs := template.FuncMap{
		"comment": comment,
		"ipport":  ipport,
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
	tmpls, err = tmpls.Funcs(funcs).ParseFS(tpl, "templates/*.cfg")
	if err != nil {
		return
	}

	err = addExtraTemplates(tmpls, "templates")
	if err != nil {
		return
	}

	if false {
		repr.Println(tmpls)
	}

	t := tmpls.Lookup("main.cfg")
	err = t.Execute(os.Stdout, struct {
		Cfg f5.F5Config
	}{
		Cfg: cfg,
	})
	repr.Println(t.DefinedTemplates())
	return
}
