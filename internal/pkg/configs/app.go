package configs

import (
	"github.com/pkg/errors"
)

type app struct {
	AgentID     string `yaml:"agent_id"`
	Environment string `yaml:"environment"`
	DebugLog    bool   `yaml:"debug_log"`
}

func (a *app) initiate() error {
	return a.validate()
}

func (a *app) validate() error {
	if a.Environment == "" {
		return errors.New("empty `environment`")
	}

	if a.AgentID == "" {
		return errors.New("empty agent_id")
	}

	return nil
}
