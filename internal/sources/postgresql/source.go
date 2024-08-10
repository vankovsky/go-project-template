package postgresql

import (
	"context"
	"time"

	"go-project-template/internal/pkg/configs"
	"go-project-template/internal/sources/log"
	"go-project-template/internal/sources/postgresql/migrations"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"
)

type Source struct {
	l            log.Logger
	pool         *pgxpool.Pool
	agentID      string
	queryTimeout time.Duration
}

func Init(cfg *configs.Root, l log.Logger) (*Source, error) {
	s := &Source{
		l:            l,
		queryTimeout: cfg.PostgreSQL.QueryTimeout,
		agentID:      cfg.App.AgentID,
	}

	connConfig, err := pgxpool.ParseConfig(cfg.PostgreSQL.URI)
	if err != nil {
		return nil, errors.Wrap(err, "parse uri")
	}

	connConfig.HealthCheckPeriod = time.Second * 5

	s.pool, err = pgxpool.NewWithConfig(context.Background(), connConfig)
	if err != nil {
		return nil, errors.Wrap(err, "connect")
	}

	if errM := s.migrate(); errM != nil {
		return nil, errors.Wrap(errM, "migrations")
	}

	return s, nil
}

func (s *Source) Close() {
	s.pool.Close()
}

func (s *Source) Acquire(ctx context.Context) (*pgxpool.Conn, error) {
	return s.pool.Acquire(ctx)
}

func (s *Source) execQuery(query string) (pgconn.CommandTag, error) {
	conn, err := s.pool.Acquire(context.Background())
	if err != nil {
		return pgconn.CommandTag{}, errors.Wrap(err, "acquire connection")
	}
	defer conn.Release()

	ctx, cancel := context.WithTimeout(context.Background(), s.queryTimeout)
	defer cancel()

	return conn.Exec(ctx, query)
}

func (s *Source) migrate() error {
	return nil
	// return s.applyMigrations(
	// 	migrations.AddAgent,
	// 	migrations.AddHost,
	// 	migrations.AddURL,
	// 	migrations.AddURLAgent,
	// 	migrations.AddScheduler,
	// 	migrations.AddCheckDNS,
	// )
}

func (s *Source) applyMigrations(migrs ...migrations.Migration) error {
	conn, err := s.pool.Acquire(context.Background())
	if err != nil {
		return errors.Wrap(err, "acquire connection")
	}
	defer conn.Release()

	var migr migrations.Migration
	for i := range migrs {
		migr = migrs[i]
		s.l.Info("migration: %v", migr.Name)
		tx, err := conn.Begin(context.Background())
		if err != nil {
			return errors.Wrap(err, "begin")
		}

		for _, statm := range migr.Statements {
			_, err = tx.Exec(context.Background(), statm)
			if err != nil {
				return errors.Wrap(err, statm)
			}
		}

		if errCommit := tx.Commit(context.Background()); errCommit != nil {
			return errors.Wrap(err, "commit")
		}
	}

	return nil
}
