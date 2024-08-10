package configs

import (
	"time"
)

const (
	defaultCheckPingInterval = 60
	defaultCheckPingWorkers  = 10
	defaultCheckPingAttempts = 15
	defaultCheckPingMaxRTT   = 10
)

type checkPing struct {
	CheckInterval time.Duration `yaml:"check_interval"`
	Workers       uint8         `yaml:"workers"`
	Attempts      uint8         `yaml:"attempts"`
	MaxRTT        time.Duration `yaml:"max_rtt"`
}

func (d *checkPing) initiate() error {
	if d.Workers == 0 {
		d.Workers = defaultCheckPingWorkers
	}

	if d.Attempts == 0 {
		d.Attempts = defaultCheckPingAttempts
	}

	if d.MaxRTT == 0 {
		d.MaxRTT = defaultCheckPingMaxRTT
	}
	d.MaxRTT *= time.Second

	if d.CheckInterval == 0 {
		d.CheckInterval = defaultCheckPingInterval * time.Second
	}

	return nil
}
