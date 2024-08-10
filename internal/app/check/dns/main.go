package dns

import (
	"context"
	"go-project-template/internal/models/log"
	"go-project-template/internal/pkg/configs"
	"go-project-template/internal/sources/postgresql"
	"net"
	"reflect"
	"sync"
	"time"

	"github.com/miekg/dns"
	"github.com/pkg/errors"
)

type DNSChecker interface {
	CheckDNSHostsToCheck() (postgresql.Hosts, error)
	CheckDNSSaveHostIPs(records postgresql.DNSResponses) error
}

type app struct {
	cfg     *configs.Root
	pg      DNSChecker
	l       log.Logger
	chHosts chan postgresql.Host
}

// Run launches workers which are retrieve IPs for the domain name and monitor
// its changes.
func Run(ctx context.Context, cfg *configs.Root, pg DNSChecker, l log.Logger) {
	a := &app{
		cfg:     cfg,
		pg:      pg,
		l:       l,
		chHosts: make(chan postgresql.Host),
	}

	a.l.Info("start")
	defer a.l.Info("exit")

	var wg sync.WaitGroup
	wg.Add(int(a.cfg.Checks.DNS.Workers))
	for i := uint8(0); i < a.cfg.Checks.DNS.Workers; i++ {
		go a.checker(&wg)
	}

	a.hostsToCheck(ctx)
	wg.Wait()
}

func (a *app) hostsToCheck(ctx context.Context) {
	a.l.Info("hosts to check: started")
	defer a.l.Info("hosts to check: exited")

	a.l.Info("interval %v", a.cfg.Checks.DNS.CheckInterval)
	for {
		select {
		case <-time.After(a.cfg.Checks.DNS.CheckInterval):
			// get hosts
			hosts, err := a.pg.CheckDNSHostsToCheck()
			if err != nil {
				a.l.Error("get hosts %v", err)
				continue
			}

			for _, h := range hosts {
				a.l.Info("check host %+v", h)
				a.chHosts <- *h
			}

		case <-ctx.Done():
			close(a.chHosts)
			return
		}
	}
}

func (a *app) checker(wg *sync.WaitGroup) {
	defer wg.Done()

	a.l.Info("checker: started")
	defer a.l.Info("checker: exit")

	for host := range a.chHosts {
		responses := a.processHost(&host)
		if err := a.pg.CheckDNSSaveHostIPs(responses); err != nil {
			a.l.Error("save result: %v", err)
		}
	}
}

func (a *app) processHost(host *postgresql.Host) postgresql.DNSResponses {
	ch := make(chan postgresql.DNSResponse, len(host.DNSServers))
	var wg sync.WaitGroup
	for _, dnsServer := range host.DNSServers {
		wg.Add(1)
		a.l.Info("[%v] resolving against %v", host.Name, dnsServer)
		go a.resolve(&wg, host.Name, dnsServer, ch)
	}

	wg.Wait()
	close(ch)

	var result postgresql.DNSResponses
	for r := range ch {
		result = append(result, r)
	}
	return result
}

func (a *app) resolve(
	wg *sync.WaitGroup,
	hostName string,
	dnsServer net.IP,
	resultCh chan<- postgresql.DNSResponse,
) {
	defer wg.Done()

	result := postgresql.DNSResponse{
		Host:      hostName,
		DNSServer: dnsServer,
	}

	config := dns.ClientConfig{
		Servers:  []string{dnsServer.String()},
		Port:     "53",
		Ndots:    1,
		Timeout:  10,
		Attempts: 3,
	}

	c := new(dns.Client)
	m := new(dns.Msg)

	m.SetQuestion(dns.Fqdn(hostName), dns.TypeA)

	r, _, err := c.Exchange(m, config.Servers[0]+":"+config.Port)
	if err != nil {
		a.l.Error("%v: exchange %v", hostName, err)
		result.Error = errors.Wrap(err, "exchange")
		resultCh <- result
		return
	}

	if r.Rcode != dns.RcodeSuccess {
		a.l.Error("BUG: non-equal answer: got %v, must: %v", r.Rcode, dns.RcodeSuccess)
		result.Error = errors.Errorf("BUG: non-equal answer: got %v, must: %v", r.Rcode, dns.RcodeSuccess)
		resultCh <- result
		return
	}

	for _, k := range r.Answer {
		switch recType := k.(type) {
		case *dns.A:
			result.Records = append(result.Records, postgresql.DNSRecord{
				IP:   recType.A,
				Type: 1, // see the types in wikipedia or in the site project
			})
		default:
			a.l.Error("%v: unexpected A-type: %v", hostName, reflect.TypeOf(k))
		}
	}

	resultCh <- result
}
