package haproxy

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"

	"github.com/Masterminds/sprig"
)

func add(a, b int) int {
	return a + b
}
func sub(a, b int) int {
	return a - b
}

func comment(str string, indent int, comment string) string {
	return spacedComment(str, 1, indent, comment)
}

func spacedComment(str string, space, indent int, comment string) string {
	idstr := strings.Repeat(" ", indent)
	spstr := strings.Repeat(" ", space)
	lines := strings.Split(str, "\n")
	for i, line := range lines {
		lines[i] = idstr + comment + spstr + line
	}
	return strings.Join(lines, "\n")
}

func indent(str string, indent int) string {
	idstr := strings.Repeat(" ", indent)
	lines := strings.Split(str, "\n")
	for i, line := range lines {
		lines[i] = idstr + line
	}
	return strings.Join(lines, "\n")
}

func split(str, sep string) []string {
	strs := strings.Split(str, sep)
	return strs
}

func stripport(str string) string {
	strs := strings.Split(str, ":")
	return strs[0]
}

func ipport(str string) string {
	strs := strings.Split(str, "/")
	ipport := strings.Split(strs[len(strs)-1], ":")
	ipport[0] = strings.Split(ipport[0], "%")[0]

	return strings.Join(ipport, ":")
}

func normalize(str string) string {
	if str[0] == '/' {
		str = str[1:]
	}
	str = strings.Replace(str, "/", "::", -1)
	return str

}

func loadTemplates(config *Config) (tmpls *template.Template, err error) {
	tmpls = template.New("")
	funcs := template.FuncMap{
		"add":       add,
		"sub":       sub,
		"comment":   comment,
		"scomment":  spacedComment,
		"indent":    indent,
		"stripport": stripport,
		"split":     split,
		"ipport":    ipport,
		"normalize": normalize,
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
		"templateIndent": func(level int64, name string, pipeline interface{}) (string, error) {
			if !config.ExpandTemplates && level != 0 {
				return fmt.Sprintf(`{{   templateIndent %d "%s" "%s" }}`, level, name, pipeline), nil
			} // else {
			// 	fmt.Printf("expending %s: %t && %t, %i\n", name, !config.ExpandTemplates, level != 0, level)
			// }

			lvl := int(level)
			if lvl < 0 {
				lvl = 0
			}

			// Disabled template
			switch p := pipeline.(type) {
			case string:
				if strings.HasPrefix(p, "disabled:") {
					idstr := strings.Repeat(" ", lvl)
					return fmt.Sprintf("%s# %s | %s", idstr, name, p), nil
				}
			}

			t := tmpls.Lookup(name)
			if t == nil {
				return fmt.Sprintf("%s### Template %s not found", strings.Repeat(" ", lvl), name), err
			}

			buf := &bytes.Buffer{}
			err := t.Execute(buf, pipeline)
			if err != nil {
				return "", err
			}

			idstr := strings.Repeat(" ", lvl)
			lines := strings.Split(buf.String(), "\n")
			for i, line := range lines {
				lines[i] = idstr + line
			}
			return strings.Join(lines, "\n"), nil
		},
	}
	tmpls = tmpls.Funcs(sprig.TxtFuncMap())
	tmpls = tmpls.Funcs(funcs)
	tmpls, err = tmpls.ParseFS(tpl, "templates/*.tpl.cfg")
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
