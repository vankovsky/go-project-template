package configs

import "github.com/pkg/errors"

type checks struct {
	DNS  checkDNS  `yaml:"dns"`
	Ping checkPing `yaml:"ping"`
	HTTP checkHTTP `yaml:"http"`
}

func (c *checks) initiate() error {
	if err := c.DNS.initiate(); err != nil {
		return errors.Wrap(err, "dns")
	}

	if err := c.Ping.initiate(); err != nil {
		return errors.Wrap(err, "ping")
	}

	if err := c.HTTP.initiate(); err != nil {
		return errors.Wrap(err, "http")
	}

	return nil
}
