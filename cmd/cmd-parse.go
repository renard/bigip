package cmd

import (
	"bigip/f5"
	"bigip/haproxy"
	"bigip/internal/log"
)

type Parse struct {
	Templates    []string `short:"t" help:"Custom template directory" type:"string"`
	GenTemplates bool     `short:"g" help:"Generate templates" type:"bool"`
	Files        []string `arg help:"Configuration files" type:"string"`
}

func (c *Parse) Run(clictx *CLIContext) (err error) {
	log.Debug("Parsing configuration files %#v", c.Files)

	hap := &haproxy.Config{
		TemplateDir: c.Templates,
	}

	cfg, err := f5.ParseFile(c.Files[0])
	if err != nil {
		return
	}
	if c.GenTemplates {
		err = haproxy.GenerateTemplates(hap, cfg)
		return
	}
	err = haproxy.Render(hap, cfg)
	return
}
