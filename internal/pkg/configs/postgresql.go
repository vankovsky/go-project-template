package configs

import (
	"time"

	"github.com/pkg/errors"
)

const defaultQueryTimeout = 60

type postgreSQL struct {
	URI          string        `yaml:"uri"`
	QueryTimeout time.Duration `yaml:"query_timeout"`
}

func (p *postgreSQL) initiate() error {
	if err := p.validate(); err != nil {
		return err
	}

	if p.QueryTimeout == 0 {
		p.QueryTimeout = defaultQueryTimeout
	}

	p.QueryTimeout *= time.Second

	return nil
}

func (p *postgreSQL) validate() error {
	if p.URI == "" {
		return errors.New("empty uri")
	}

	if p.QueryTimeout < 0 {
		return errors.New("negative `query_timeout_sec`")
	}

	return nil
}
