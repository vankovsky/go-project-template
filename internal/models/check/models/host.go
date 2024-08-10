package models

import (
	"net"
	"time"

	"github.com/google/uuid"
)

type DNSRecord struct {
	Type string
	IP   net.IP
}

type DNSResponse struct {
	DNSServer net.IP
	Records   []*DNSRecord
	Error     error
}

type Host struct {
	ID           uuid.UUID
	Name         string
	DNSResponses []*DNSResponse
}

type PingResult struct {
	IP           net.IP
	ResponseTime time.Duration
}

type PingedIPs struct {
	ID    int
	IPs   []PingResult
	Error error
}

type URL struct {
	ID           uuid.UUID
	URL          string
	Error        error
	StatusCode   uint8
	ResponseTime uint16
	ResponseSize uint16
}
