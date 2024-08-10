package configs

import (
	"time"

	"github.com/pkg/errors"
)

const (
	defaultCheckHTTPInterval = 1
	defaultCheckHTTPWorkers  = 1
)

type checkHTTP struct {
	CheckInterval time.Duration `yaml:"check_interval"`
	Workers       uint8         `yaml:"workers"`
}

func (d *checkHTTP) initiate() error {
	if err := d.validate(); err != nil {
		return err
	}

	if d.Workers == 0 {
		d.Workers = defaultCheckPingWorkers
	}

	if d.CheckInterval == 0 {
		d.CheckInterval = time.Duration(defaultCheckHTTPInterval) * time.Second
	}

	return nil
}

func (d *checkHTTP) validate() error {
	if d.CheckInterval < 0 {
		return errors.New("negative check_interval_sec")
	}

	if d.Workers < 0 {
		return errors.New("negative workers")
	}

	return nil
}
