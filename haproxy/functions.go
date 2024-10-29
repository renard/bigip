// Copyright © 2023 Sébastien Gross <seb•ɑƬ•chezwam•ɖɵʈ•org>
//
// Created: 2021-12-31
// Last changed: 2024-10-29 11:49:09
//
// This program is free software: you can redistribute it and/or
// modify it under the terms of the GNU Affero General Public License
// as published by the Free Software Foundation, either version 3 of
// the License, or (at your option) any later version.
//
// This program is distributed in the hope that it will be useful, but
// WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU
// Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public
// License along with this program. If not, see
// <http://www.gnu.org/licenses/>.

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
	str = strings.Replace(str, " ", "_", -1)
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
				config.Log.Error("Template %s not found", name)
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
