// Copyright © 2023 Sébastien Gross <seb•ɑƬ•chezwam•ɖɵʈ•org>
//
// Created: 2021-12-19
// Last changed: 2024-10-29 11:47:52
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

package cmd

import (
	"bigip/f5"
	"bigip/haproxy"
)

type Convert struct {
	Templates []string `short:"t" help:"Custom template directory" type:"string"`
	Files     []string `arg help:"Configuration files" type:"string"`
	OutputDir string   `short:"o" help:"Output directory" default:"" type:"string"`
	Virtual   []string `short:"V" help:"Limit conversion to specific virtuals" type:"string"`
	Pool      []string `short:"P" help:"Add extra pools when virtual flag is used" type:"string"`
}

func (c *Convert) Run(clictx *CLIContext) (err error) {
	hap := &haproxy.Config{
		Files:           c.Files,
		TemplateDir:     c.Templates,
		ExpandTemplates: true,
		OutputDir:       c.OutputDir,
		Virtual:         c.Virtual,
		Pool:            c.Pool,
		Log:             clictx.log,
	}

	cfg, err := f5.ParseFile(clictx.log, c.Files)
	if err != nil {
		return
	}
	err = haproxy.Render(hap, cfg)
	return
}
