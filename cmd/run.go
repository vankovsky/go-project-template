package main

import (
	"context"
	"flag"
	"fmt"
	"go-project-template/internal/app/check/dns"
	"go-project-template/internal/pkg/configs"
	"go-project-template/internal/sources"
	"go-project-template/internal/sources/log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
)

type mainApp struct {
	cfg *configs.Root
	src *sources.Sources
	l   log.Logger
}

func newMainApp(configFilename string) (*mainApp, error) {
	// initiate settings file:
	cfg, err := configs.Init(configFilename)
	if err != nil {
		return nil, errors.Wrap(err, "config error")
	}

	return &mainApp{cfg: cfg}, nil
}

func (a *mainApp) initSources() error {
	src, err := sources.Init(a.cfg)
	if err != nil {
		return errors.Wrap(err, "sources init")
	}

	a.src = src
	a.l = log.New("main")

	return nil
}

func (a *mainApp) runApps(ctx context.Context) {
	g, _ := errgroup.WithContext(ctx)
	g.Go(func() error { dns.Run(ctx, a.cfg, a.src.PG, log.New("check/dns")); return nil })
	// g.Go(func() error { ping.Run(ctx, a.cfg, a.src.PG, log.New("check/ping")); return nil })
	// g.Go(func() error { http.Run(ctx, a.cfg, a.src.PG, log.New("check/http")); return nil })
	// g.Go(func() error { return status.Run(ctx, a.cfg, a.src, log.New("status")) })

	if err := g.Wait(); err != nil {
		a.l.Error("abnormal exit: %v", err)
	}
}

func main() {
	// CLI flags:
	var configFile string
	var doValidation bool
	var doMigrations bool

	flag.StringVar(&configFile, "config", "", "configuration file path")
	flag.BoolVar(&doValidation, "validate", false, "validate configuration and exit")
	flag.BoolVar(&doMigrations, "migrations", false, "apply DB's migrations and exit")

	flag.Parse()

	a, err := newMainApp(configFile)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	if doValidation {
		os.Exit(0)
	}

	// initiate connections to DBs:
	if err := a.initSources(); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	if doMigrations {
		os.Exit(0)
	}

	// define the context and Ctrl-C exit:
	ctx, cancel := context.WithCancel(context.Background())
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGTERM, syscall.SIGINT)

	// cleaning-up the listeners:
	defer func() {
		signal.Stop(c)
		cancel()
	}()

	// catch the interruptions:
	go func() {
		<-c
		a.l.Info("got cancel signal")
		cancel()
	}()

	// catch panics:
	if a.cfg.Sentry.DSN != "" {
		defer func() {
			if err := recover(); err != nil {
				sentry.CurrentHub().Recover(err)
				sentry.Flush(time.Second * 5)
			}
		}()
	}

	a.runApps(ctx)
}
