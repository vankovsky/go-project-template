package pinger

import "time"

type Hosts struct {
	UpdatedAt time.Time
	Hosts     []string
}

func (ph *Hosts) Update(hosts []string) {
	if len(hosts) == 0 {
		return
	}

	for _, h := range hosts {
		if ph.IsHostExisted(h) {
			continue
		}

		ph.Hosts = append(ph.Hosts, h)
	}
}

func (ph *Hosts) IsHostExisted(host string) bool {
	for _, h := range ph.Hosts {
		if h == host {
			return true
		}
	}

	return false
}

func (ph *Hosts) IsEmpty() bool {
	if ph == nil {
		return true
	}

	return len(ph.Hosts) == 0
}
