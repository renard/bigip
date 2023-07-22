// Copyright © 2023 Sébastien Gross <seb•ɑƬ•chezwam•ɖɵʈ•org>
//
// Created: 2021-12-19
// Last changed: 2023-07-22 02:59:24
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
	"embed"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"bigip/f5"

	"github.com/alecthomas/repr"
	"github.com/jackdoe/go-pager"
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
	Virtual         []string
	Pool            []string
}

func addExtraTemplates(t *template.Template, dir string) (err error) {
	err = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if strings.HasSuffix(path, ".tpl.cfg") {
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
	var out io.Writer
	// Multiple ways to output converted file
	switch config.OutputDir {
	case "":
		// pager
		close := func() {}
		out, close = pager.Pager("less", "more", "cat")
		defer close()
	case "-":
		// stdout
		out = os.Stdout
	default:
		// file
		out, err = os.Create(fmt.Sprintf("%s/main.cfg", config.OutputDir))
		defer out.(*os.File).Close()
	}

	tmpls, err := loadTemplates(config)
	if err != nil {
		return
	}

	if false {
		repr.Println(tmpls)
	}

	t := tmpls.Lookup("main")

	f5c := f5.NewF5Config()
	if len(config.Virtual) == 0 {
		f5c.LtmVirtual = f5config.LtmVirtual
		f5c.LtmPool = f5config.LtmPool
	} else {
		// Limit to provided virtual only
		for _, virtual := range config.Virtual {
			if v, ok := f5config.LtmVirtual[virtual]; ok {
				f5c.LtmVirtual[virtual] = v
				defaultPool := v.(*f5.LtmVirtual).Pool

				// Only keep pools used by selected virtual
				if p, pok := f5config.LtmPool[defaultPool]; pok {
					f5c.LtmPool[defaultPool] = p
				}

				// TODO: Add extra pools.
				for _, pool := range config.Pool {
					if p, pok := f5config.LtmPool[pool]; pok {
						f5c.LtmPool[pool] = p
					}
				}
			}
		}
	}
	f5c.LtmNode = f5config.LtmNode
	f5c.LtmRule = f5config.LtmRule
	f5c.LtmProfile = f5config.LtmProfile
	f5c.LtmMonitor = f5config.LtmMonitor
	f5c.LtmPersistence = f5config.LtmPersistence

	err = t.Execute(out, struct {
		F5config f5.F5Config
		Config   Config
	}{
		F5config: f5c,
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

	exports := config.Export
	if len(exports) == 0 {
		exports = []string{"virtual", "pool", "rule", "policy", "profile", "node", "monitor",
			"persistence"}
	}
	t := tmpls.Lookup("export")
	for _, tp := range exports {
		f5c := f5.NewF5Config()
		switch tp {
		case "virtual":
			f5c.LtmVirtual = f5config.LtmVirtual
		case "pool":
			f5c.LtmPool = f5config.LtmPool
			f5c.LtmNode = f5config.LtmNode
		case "rule":
			f5c.LtmRule = f5config.LtmRule
		case "policy":
			f5c.LtmPolicy = f5config.LtmPolicy
		case "profile":
			f5c.LtmProfile = f5config.LtmProfile
		case "node":
			f5c.LtmNode = f5config.LtmNode
		case "monitor":
			f5c.LtmMonitor = f5config.LtmMonitor
		case "persistence":
			f5c.LtmPersistence = f5config.LtmPersistence
		}

		fh, e := os.Create(fmt.Sprintf("%s/%s.tpl.cfg", config.OutputDir, tp))
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
