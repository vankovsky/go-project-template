package ping

import (
	"context"
	"net"
	"sync"
	"time"

	"go-project-template/internal/models/log"
	"go-project-template/internal/pkg/configs"
	"go-project-template/internal/sources/postgresql"

	"github.com/go-ping/ping"
)

type app struct {
	cfg *configs.Root
	pg  postgresql.PingChecker
	l   log.Logger
}

func Run(ctx context.Context, cfg *configs.Root, pg postgresql.PingChecker, l log.Logger) {
	a := &app{
		cfg: cfg,
		pg:  pg,
		l:   l,
	}

	a.l.Info("start")
	defer a.l.Info("exit")

	ch := make(chan hostIP, a.cfg.Checks.Ping.Workers)
	wg := sync.WaitGroup{}
	for i := uint8(0); i < a.cfg.Checks.Ping.Workers; i++ {
		wg.Add(1)
		go a.worker(&wg, ch)
	}

	a.hostRetrieve(ctx, ch)
	close(ch)
	wg.Wait()
}

func (a *app) hostRetrieve(ctx context.Context, ch chan<- hostIP) {
	ticker := time.NewTicker(a.cfg.Checks.Ping.CheckInterval)
	for {
		select {
		case <-ticker.C:
			// get hosts
			hosts, err := a.pg.CheckPingHosts()
			if err != nil {
				a.l.Error("get hosts: %v", err)
				continue
			}
			a.l.Info("hosts to check %d", len(hosts))

			var wgProc sync.WaitGroup
			for _, h := range hosts {
				for _, ip := range h.IPs {
					ch <- hostIP{hostID: h.HostID, ip: ip}
				}
			}
			wgProc.Wait()

		case <-ctx.Done():
			return
		}
	}
}

type hostIP struct {
	hostID string
	ip     net.IP
}

func (a *app) worker(wg *sync.WaitGroup, ch <-chan hostIP) {
	defer wg.Done()

	for h := range ch {
		a.processIP(h)
	}
}

func (a *app) processIP(h hostIP) {
	result := postgresql.PingResult{
		HostID: h.hostID,
		IP:     h.ip,
	}

	pinger, err := ping.NewPinger(result.IP.String())
	if err != nil {
		a.l.Error("new pinger: %v", err)
		return
	}

	pinger.SetPrivileged(true)
	pinger.Count = 3
	pinger.Timeout = time.Duration(pinger.Count) * a.cfg.Checks.Ping.MaxRTT

	pinger.OnFinish = func(stats *ping.Statistics) {
		result.Stats = stats
		if err := a.pg.CheckPingSaveIPResult(&result); err != nil {
			a.l.Error("save result: %v", err)
			return
		}
	}

	if err := pinger.Run(); err != nil {
		a.l.Error("run pinger: %v", err)
		return
	}
}
