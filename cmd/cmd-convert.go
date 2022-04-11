package cmd

import (
	"bigip/f5"
	"bigip/haproxy"
	"bigip/internal/log"
)

type Convert struct {
	Templates []string `short:"t" help:"Custom template directory" type:"string"`
	Files     []string `arg help:"Configuration files" type:"string"`
	OutputDir string   `short:"o" help:"Output directory" default:"-" type:"string"`
	Virtual   []string `short:"V" help:"Limit conversion to specific virtuals" type:"string"`
}

func (c *Convert) Run(clictx *CLIContext) (err error) {
	log.Debug("Parsing configuration files %#v", c.Files)

	hap := &haproxy.Config{
		Files:           c.Files,
		TemplateDir:     c.Templates,
		ExpandTemplates: true,
		OutputDir:       c.OutputDir,
		Virtual:         c.Virtual,
	}

	cfg, err := f5.ParseFile(c.Files)
	if err != nil {
		return
	}
	err = haproxy.Render(hap, cfg)
	return
}
