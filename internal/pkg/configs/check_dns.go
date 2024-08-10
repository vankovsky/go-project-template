package configs

import (
	"net"
	"time"

	"github.com/pkg/errors"
)

const (
	defaultCheckDNSInterval = 60 * time.Second
	defaultCheckDNSWorkers  = 1
)

type checkDNS struct {
	DNSHosts      []net.IP      `yaml:"dns_hosts"`
	CheckInterval time.Duration `yaml:"check_interval"`
	Workers       uint8         `yaml:"workers"`
}

func (d *checkDNS) initiate() error {
	if err := d.validate(); err != nil {
		return err
	}

	if d.Workers == 0 {
		d.Workers = defaultCheckDNSWorkers
	}

	if d.CheckInterval == 0 {
		d.CheckInterval = defaultCheckDNSInterval
	}

	return nil
}

func (d *checkDNS) validate() error {
	if d.CheckInterval < 0 {
		return errors.New("negative check_schedule_interval_sec")
	}

	// if len(d.DNSHosts) == 0 {
	// 	return errors.New("empty dns_hosts")
	// }

	if d.Workers < 0 {
		return errors.New("negative workers")
	}

	return nil
}
