// Copyright © 2023 Sébastien Gross <seb•ɑƬ•chezwam•ɖɵʈ•org>
//
// Created: 2021-12-19
// Last changed: 2023-07-22 02:57:08
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
	"bigip/internal/log"

	"github.com/alecthomas/kong"
)

type CLIContext struct {
	//cli     string
	Verbose int `help:"Run in verbose mode." short:"v" type:"counter"`
}

var CLI struct {
	CLIContext
	Convert      Convert      `cmd help:"Convert F5 configuration"`
	GenTemplates GenTemplates `cmd help:"Generate block configuration templates"`
	// Merge   Merge   `cmd help:"Merge SQLite results"`
	// Graph   Graph   `cmd help:"Graph bench results"`
	// LVS     LVS     `cmd help:"LVS tests"`
	// Version Version `cmd help:"Show version"`
	// TestYAML TestYAML `cmd help:"Test configuration"`
}

func ParseCli() {
	ctx := kong.Parse(&CLI,
		kong.UsageOnError(),
	)
	log.SetLevel(CLI.Verbose)
	err := ctx.Run(&CLIContext{
		Verbose: CLI.Verbose,
		//cli:     tools.QuoteShell(os.Args, nil),
	})
	ctx.FatalIfErrorf(err)
}
