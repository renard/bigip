package cmd

import (
	"bigip/f5"
	"bigip/haproxy"
	"bigip/internal/log"
)

type Convert struct {
	Templates []string `short:"t" help:"Custom template directory" type:"string"`
	Files     []string `arg help:"Configuration files" type:"string"`
}

func (c *Convert) Run(clictx *CLIContext) (err error) {
	log.Debug("Parsing configuration files %#v", c.Files)

	hap := &haproxy.Config{
		Files:           c.Files,
		TemplateDir:     c.Templates,
		ExpandTemplates: true,
	}

	cfg, err := f5.ParseFile(c.Files)
	if err != nil {
		return
	}
	err = haproxy.Render(hap, cfg)
	return
}
