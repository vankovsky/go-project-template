package configs

import (
	"github.com/pkg/errors"
)

const (
	PrometheusPrefix = "d"
)

// Root is a type for the file configuration format.
type Root struct {
	App        app        `yaml:"app"`
	HttpServer httpServer `yaml:"http_server"`
	Sentry     sentry     `yaml:"sentry"`
	PostgreSQL postgreSQL `yaml:"postgresql"`
	Checks     checks     `yaml:"checks"`
}

func (r *Root) initiate() error {
	if err := r.App.initiate(); err != nil {
		return errors.Wrap(err, "app")
	}

	if err := r.HttpServer.validate(); err != nil {
		return errors.Wrap(err, "http server")
	}

	if err := r.PostgreSQL.initiate(); err != nil {
		return errors.Wrap(err, "postgresql")
	}

	if err := r.Checks.initiate(); err != nil {
		return errors.Wrap(err, "checks")
	}

	return nil
}
