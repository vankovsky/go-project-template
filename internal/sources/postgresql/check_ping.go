package postgresql

import (
	"context"
	"net"
	"sort"
	"time"

	"github.com/go-ping/ping"
	"github.com/jackc/pgx/v5"
	"github.com/pkg/errors"
)

type PingChecker interface {
	CheckPingHosts() ([]*PingHost, error)
	CheckPingSaveIPResult(result *PingResult) error
}

type PingHost struct {
	HostID string
	IPs    []net.IP
}

type PingResult struct {
	Stats  *ping.Statistics
	HostID string
	IP     net.IP
}

func (p *PingResult) calculate95thPercentile() time.Duration {
	if len(p.Stats.Rtts) == 0 {
		return 0
	}

	sort.Slice(p.Stats.Rtts, func(i, j int) bool {
		return p.Stats.Rtts[i] < p.Stats.Rtts[j]
	})
	index := int(float64(len(p.Stats.Rtts)) * 0.95)
	return p.Stats.Rtts[index]
}

func (p *PingResult) calculate50thPercentile() time.Duration {
	if len(p.Stats.Rtts) == 0 {
		return 0
	}
	sort.Slice(p.Stats.Rtts, func(i, j int) bool {
		return p.Stats.Rtts[i] < p.Stats.Rtts[j]
	})
	index := int(float64(len(p.Stats.Rtts)) * 0.5)
	return p.Stats.Rtts[index]
}

func (p *PingResult) toMicroseconds() []int {
	res := make([]int, len(p.Stats.Rtts))
	for i, d := range p.Stats.Rtts {
		res[i] = int(d / time.Microsecond)
	}
	return res
}

func (s *Source) CheckPingHosts() ([]*PingHost, error) {
	conn, err := s.pool.Acquire(context.Background())
	if err != nil {
		return nil, errors.Wrap(err, "acquire connection")
	}
	defer conn.Release()

	ctx, cancel := context.WithTimeout(context.Background(), s.queryTimeout)
	defer cancel()

	rows, err := conn.Query(ctx, `SELECT host_id, array_agg(ip) FROM host_hostip group by host_id`)
	if err != nil {
		return nil, errors.Wrap(err, "execute")
	}

	var hosts []*PingHost
	for rows.Next() {
		h := PingHost{}
		err := rows.Scan(&h.HostID, &h.IPs)
		if err != nil {
			return nil, errors.Wrap(err, "scan result")
		}
		hosts = append(hosts, &h)
	}

	return hosts, nil
}

func (s *Source) CheckPingSaveIPResult(r *PingResult) error {
	conn, err := s.pool.Acquire(context.Background())
	if err != nil {
		return errors.Wrap(err, "acquire connection")
	}
	defer conn.Release()

	ctx, cancel := context.WithTimeout(context.Background(), s.queryTimeout)
	defer cancel()

	err = pgx.BeginFunc(ctx, conn, func(tx pgx.Tx) error {
		if _, err := tx.Exec(ctx, `
		insert into check_ping (host_id, agent_id, ip, "avg_value", p50, p95, lost_packets, "min_value", "max_value", stddev_value, duplicated_packets, sent_packets) VALUES 
	($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
on CONFLICT (host_id, agent_id, ip) do update set updated_at = now(), 
"avg_value" = $4, p50 = $5, p95 = $6, lost_packets = $7, "min_value" = $8, "max_value" = $9, stddev_value = $10, duplicated_packets = $11, sent_packets = $12
		`,
			r.HostID,
			s.agentID,
			r.IP.String(),
			r.Stats.AvgRtt.Microseconds(),
			r.calculate50thPercentile().Microseconds(),
			r.calculate95thPercentile().Microseconds(),
			int(r.Stats.PacketLoss*100),
			r.Stats.MinRtt.Microseconds(),
			r.Stats.MaxRtt.Microseconds(),
			r.Stats.StdDevRtt.Microseconds(),
			r.Stats.PacketsRecvDuplicates,
			len(r.Stats.Rtts),
		); err != nil {
			return err
		}

		// check_ping_host_id_7d5dff5f_fk_host_url_id

		if _, err := tx.Exec(ctx, `
		insert into check_pinglog (host_id, agent_id, ip, "avg_value", p50, p95, lost_packets, "min_value", "max_value", stddev_value, duplicated_packets, rtts) VALUES 
	($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
		`,
			r.HostID,
			s.agentID,
			r.IP.String(),
			r.Stats.AvgRtt.Microseconds(),
			r.calculate50thPercentile().Microseconds(),
			r.calculate95thPercentile().Microseconds(),
			int(r.Stats.PacketLoss*100),
			r.Stats.MinRtt.Microseconds(),
			r.Stats.MaxRtt.Microseconds(),
			r.Stats.StdDevRtt.Microseconds(),
			r.Stats.PacketsRecvDuplicates,
			r.toMicroseconds(),
		); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return errors.Wrap(err, "tx")
	}

	return nil
}
