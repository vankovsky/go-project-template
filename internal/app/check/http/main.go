package http

import (
	"context"
	"go-project-template/internal/pkg/configs"
	"go-project-template/internal/sources/log"
	"go-project-template/internal/sources/postgresql"
	"net/http"
	"sync"
	"time"
)

type app struct {
	cfg   *configs.Root
	pg    postgresql.HTTPChecker
	l     log.Logger
	chURL chan postgresql.URL
}

func Run(ctx context.Context, cfg *configs.Root, pg postgresql.HTTPChecker, l log.Logger) {
	a := &app{
		cfg:   cfg,
		pg:    pg,
		l:     l,
		chURL: make(chan postgresql.URL),
	}

	a.l.Info("start")
	defer a.l.Info("exit")

	wg := sync.WaitGroup{}
	for i := uint8(0); i < a.cfg.Checks.HTTP.Workers; i++ {
		wg.Add(1)
		go a.worker(&wg)
	}

	a.urlsToCheck(ctx)
	close(a.chURL)
	wg.Wait()
}

func (a *app) urlsToCheck(ctx context.Context) {
	a.l.Info("urlsToCheck start")
	defer a.l.Info("urlsToCheck exit")

	ticker := time.NewTicker(a.cfg.Checks.HTTP.CheckInterval)
	for {
		select {
		case <-ticker.C:
			urls, err := a.pg.CheckHTTPURLs()
			if err != nil {
				a.l.Error("get urls: %v", err)
				continue
			}
			a.l.Info("urls to check %d", len(urls))
			for _, u := range urls {
				a.chURL <- *u
			}

		case <-ctx.Done():
			return
		}
	}
}

func (a *app) worker(wg *sync.WaitGroup) {
	defer wg.Done()

	for u := range a.chURL {
		a.processURL(u)
	}
}

func (a *app) processURL(u postgresql.URL) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(
		ctx, http.MethodGet, u.URL, nil,
	)

	c := &http.Client{}
	result := postgresql.HTTPResponse{URLID: u.ID}
	if err != nil {
		a.l.Error("new request: %v", err)
		result.Error = err
		if err := a.pg.CheckHTTPSaveResponse(&result); err != nil {
			a.l.Error("save response: %v", err)
		}
		return
	}

	t := time.Now()
	resp, err := c.Do(req)
	result.ResponseTime = time.Since(t)
	if err != nil {
		a.l.Error("do request: %v", err)
		result.Error = err
		if err := a.pg.CheckHTTPSaveResponse(&result); err != nil {
			a.l.Error("save response: %v", err)
		}
		return
	}

	result.StatusCode = resp.StatusCode

	defer func() {
		if err := resp.Body.Close(); err != nil {
			a.l.Error("close body: %v")
		}
	}()

	result.ResponseSize = len(result.Body)
	if err := a.pg.CheckHTTPSaveResponse(&result); err != nil {
		a.l.Error("save response: %v", err)
	}
	return
}
