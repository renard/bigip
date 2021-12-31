package cmd

import (
	"bigip/f5"
	"bigip/haproxy"
	"bigip/internal/log"
)

type Parse struct {
	Files []string `arg help:"Configuration files" type:"string"`
}

func (c *Parse) Run(clictx *CLIContext) (err error) {
	log.Debug("Parsing configuration files %#v", c.Files)

	hap := &haproxy.Config{}

	cfg, err := f5.ParseFile(c.Files[0])
	if err != nil {
		return
	}
	err = haproxy.Render(hap, cfg)
	return
}
