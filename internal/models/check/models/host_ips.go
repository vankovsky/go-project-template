package models

import (
	"fmt"
	"net"
	"strings"
)

type IP struct {
	IP   net.IP
	Type string
	Err  error
}

type HostIPs struct {
	ID       int
	Hostname string
	IPs      []IP
	Err      error
}

func (h *HostIPs) ToRawJSON() string {
	var ips []string
	for _, ip := range h.IPs {
		ip4 := ip.IP.To4()
		ips = append(ips, fmt.Sprintf("%q", ip4.String()))
	}

	return `'{"a":[` + strings.Join(ips, ",") + "]}'"
}

func (h *HostIPs) KeyRaw() string {
	var key string
	for i, ip := range h.IPs {
		if i != 0 {
			key += "," + ip.IP.String()
			continue
		}
		key = ip.IP.String()
	}

	return "'" + key + "'"
}
