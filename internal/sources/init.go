package sources

import (
	"go-project-template/internal/pkg/configs"
	"go-project-template/internal/sources/log"
	"go-project-template/internal/sources/metrics"
	"go-project-template/internal/sources/postgresql"

	"github.com/pkg/errors"
)

type Sources struct {
	PG *postgresql.Source
}

func Init(cfg *configs.Root) (*Sources, error) {
	var sources Sources
	var err error

	if err := log.Init(cfg); err != nil {
		return nil, errors.Wrap(err, "init log")
	}

	metrics.Init()

	sources.PG, err = postgresql.Init(cfg, log.New("pg"))
	if err != nil {
		return nil, errors.Wrap(err, "init client")
	}

	return &sources, nil
}
