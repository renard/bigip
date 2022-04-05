package cmd

import (
	"bigip/f5"
	"bigip/haproxy"
	"bigip/internal/log"
)

type GenTemplates struct {
	Templates       []string `short:"t" help:"Custom template directory" type:"string"`
	Files           []string `arg help:"Configuration files" type:"string"`
	Export          []string `short:"e" help:"Limit export to specific configuration parts" type:"string"`
	OutputDir       string   `short:"o" help:"Output directory" default:"-" type:"string"`
	ExpandTemplates bool     `help:"Expand templates" default:"false" type:"bool"`
}

func (c *GenTemplates) Run(clictx *CLIContext) (err error) {
	log.Debug("Parsing configuration files %#v", c.Files)

	hap := &haproxy.Config{
		TemplateDir:     c.Templates,
		Export:          c.Export,
		OutputDir:       c.OutputDir,
		ExpandTemplates: c.ExpandTemplates,
	}

	cfg, err := f5.ParseFile(c.Files)
	if err != nil {
		return
	}
	err = haproxy.GenerateTemplates(hap, cfg)
	return
}
