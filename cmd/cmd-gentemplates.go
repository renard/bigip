// Copyright © 2023 Sébastien Gross <seb•ɑƬ•chezwam•ɖɵʈ•org>
//
// Created: 2022-01-04
// Last changed: 2024-10-09 01:07:45
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
	//	"bigip/internal/log"
)

type GenTemplates struct {
	Templates       []string `short:"t" help:"Custom template directory" type:"string"`
	Files           []string `arg help:"Configuration files" type:"string"`
	Export          []string `short:"e" help:"Limit export to specific configuration parts" type:"string"`
	OutputDir       string   `short:"o" help:"Output directory" default:"-" type:"string"`
	ExpandTemplates bool     `help:"Expand templates" default:"false" type:"bool"`
}

func (c *GenTemplates) Run(clictx *CLIContext) (err error) {
	hap := &haproxy.Config{
		TemplateDir:     c.Templates,
		Export:          c.Export,
		OutputDir:       c.OutputDir,
		ExpandTemplates: c.ExpandTemplates,
	}

	cfg, err := f5.ParseFile(clictx.log, c.Files)
	if err != nil {
		return
	}
	err = haproxy.GenerateTemplates(hap, cfg)
	return
}
