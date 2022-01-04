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
	Parse        Parse        `cmd help:"Parse config"`
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
