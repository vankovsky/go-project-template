package postgresql

import (
	"context"
	"fmt"
	"net"
	"strings"

	"github.com/pkg/errors"
)

type Host struct {
	Name       string
	DNSServers []net.IP
}

func (h *Host) LogString() string {
	return fmt.Sprintf("%q", h.Name)
}

type Hosts []*Host

func (h Hosts) Names() string {
	names := make([]string, 0, len(h))
	for _, host := range h {
		names = append(names, host.Name)
	}
	return strings.Join(names, ", ")
}

type DNSRecord struct {
	Value string
	IP    net.IP
	Type  int
}

// sqlString return an SQL string of the values.
func (d DNSRecord) sqlValues() string {
	return ""
}

type DNSResponse struct {
	Error     error
	Host      string
	DNSServer net.IP
	Records   []DNSRecord
}

func (d DNSResponse) IsEmpty() bool {
	return len(d.Records) == 0
}

type DNSResponses []DNSResponse

type sqlQuery struct {
	stmt string
	args []interface{}
}

func (response *DNSResponse) sqlInsertCTESubQuery(agentID string) sqlQuery {
	queryTmpl := `
    inserted as (
      INSERT INTO check_dns (
        created_at,
        record_ip,
        record_value,
        record_type,
        agent_id,
        dns_server_id,
        host_id,
        error
      ) VALUES %v
      ON CONFLICT (
        record_ip,
        record_value,
        record_type,
        agent_id,
        dns_server_id,
        host_id
      ) do nothing
      RETURNING *
    )
  `

	rowsInsertedCTE := make([]string, 0, len(response.Records))
	var cnt uint = 0
	var rowsValues []interface{}
	for _, record := range response.Records {

		rowsInsertedCTE = append(rowsInsertedCTE,
			fmt.Sprintf("(now(), $%d, $%d, $%d, $%d, $%d, $%d, $%d)",
				// now() - created_at
				cnt+1, // record_ip,
				cnt+2, // record_value
				cnt+3, // record_type
				cnt+4, // agent_id
				cnt+5, // dns_server_id
				cnt+6, // host_id
				cnt+7, // error
			),
		)
		cnt += 7

		rowsValues = append(
			rowsValues,
			// now() - created_at
			record.IP,          // record_ip
			record.Value,       // record_value
			record.Type,        // record_type
			agentID,            // agent_id
			response.DNSServer, // dns_server_id
			response.Host,      // host_id
			response.Error,     // error
		)
	}
	return sqlQuery{
		stmt: fmt.Sprintf(queryTmpl, strings.Join(rowsInsertedCTE, ",")),
		args: rowsValues,
	}
}

func (response *DNSResponse) sqlDeleteCTESubQuery(agentID string, cnt int) sqlQuery {
	queryTmpl := `
  deleted as (
    DELETE FROM check_dns WHERE
        host_id = $%d and
        agent_id = $%d and
        dns_server_id = $%d and
        (record_ip, record_value, record_type) NOT IN (%%v)
    RETURNING *
  )
  `

	query := fmt.Sprintf(queryTmpl, cnt+1, cnt+2, cnt+3)
	cnt += 3

	var rowsValues []interface{}
	rowsValues = append(rowsValues, response.Host, agentID, response.DNSServer)

	whereClauseStmt := make([]string, 0, len(response.Records))
	for _, record := range response.Records {
		whereClauseStmt = append(whereClauseStmt,
			fmt.Sprintf("($%d, $%d, $%d)",
				cnt+1, // record_ip,
				cnt+2, // record_value
				cnt+3, // record_type
			),
		)
		cnt += 3

		rowsValues = append(
			rowsValues,
			record.IP,    // record_ip
			record.Value, // record_value
			record.Type,  // record_type
		)
	}
	return sqlQuery{
		stmt: fmt.Sprintf(query, strings.Join(whereClauseStmt, ",")),
		args: rowsValues,
	}
}

// sqlUpsertQuery upserts DNS-records and log changes to the log-table.
func (response *DNSResponse) sqlUpsertQuery(agentID string) sqlQuery {
	queryTmpl := `
			WITH
        %v,
        %v

			INSERT INTO check_dns_change_log (mode, record_ip, record_value, record_type, agent_id, dns_server_id, host_id)
			SELECT 2, record_ip, record_value, record_type, agent_id, dns_server_id, host_id FROM deleted
			UNION ALL
			SELECT 1, record_ip, record_value, record_type, agent_id, dns_server_id, host_id FROM inserted
 `

	insertCTE := response.sqlInsertCTESubQuery(agentID)
	deleteCTE := response.sqlDeleteCTESubQuery(agentID, len(insertCTE.args))

	return sqlQuery{
		stmt: fmt.Sprintf(queryTmpl, insertCTE.stmt, deleteCTE.stmt),
		args: append(insertCTE.args, deleteCTE.args...),
	}
}

func (response *DNSResponse) sqlUpdate(agentID string) sqlQuery {
	queryTmpl := `
			update check_dns 
      set updated_at = now()
      where
        agent_id = $1
        and host_id = $2 
        and dns_server_id = $3 
  `

	return sqlQuery{
		stmt: queryTmpl,
		args: []interface{}{agentID, response.Host, response.DNSServer},
	}
}

// CheckDNSHostsToCheck returns a list of hosts' names to get their DNS records.
func (s *Source) CheckDNSHostsToCheck() (Hosts, error) {
	conn, err := s.pool.Acquire(context.Background())
	if err != nil {
		return nil, errors.Wrap(err, "acquire connection")
	}
	defer conn.Release()

	ctx, cancel := context.WithTimeout(context.Background(), s.queryTimeout)
	defer cancel()

	rows, err := conn.Query(ctx,
		`SELECT
	h.name,
	array(SELECT ip FROM host_dnsserver) as dns_servers
FROM host_host h
INNER JOIN host_url u ON h.name = u.host_id
INNER JOIN host_url_agents ua ON ua.url_id = u.id
WHERE ua.agent_id = $1
`, s.agentID)
	if err != nil {
		return nil, errors.Wrap(err, "execute")
	}

	var hosts Hosts
	for rows.Next() {
		h := Host{}
		err := rows.Scan(&h.Name, &h.DNSServers)
		if err != nil {
			return nil, errors.Wrap(err, "scan result")
		}

		hosts = append(hosts, &h)
	}

	return hosts, nil
}

func (s *Source) CheckDNSSaveHostIPs(responses DNSResponses) error {
	if len(responses) == 0 {
		s.l.Error("BUG: empty records")
		return nil
	}

	conn, err := s.pool.Acquire(context.Background())
	if err != nil {
		return errors.Wrap(err, "acquire connection")
	}
	defer conn.Release()

	ctx, cancel := context.WithTimeout(context.Background(), s.queryTimeout)
	defer cancel()

	var query sqlQuery
	for _, r := range responses {
		query = r.sqlUpsertQuery(s.agentID)
		if _, err = conn.Exec(ctx, query.stmt, query.args...); err != nil {
			fmt.Println("stmt", query.stmt)
			fmt.Println("args", query.args)
			return errors.Wrap(err, "execute upsert")
		}

		query = r.sqlUpdate(s.agentID)
		if _, err = conn.Exec(ctx, query.stmt, query.args...); err != nil {

			fmt.Println("stmt", query.stmt)
			fmt.Println("args", query.args)
			return errors.Wrap(err, "execute update")
		}
	}

	return nil
}
