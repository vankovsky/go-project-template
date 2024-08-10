package postgresql

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/pkg/errors"
)

type HTTPChecker interface {
	CheckHTTPURLs() ([]*URL, error)
	CheckHTTPSaveResponse(result *HTTPResponse) error
}

type URL struct {
	ID  string
	URL string
}

type HTTPResponse struct {
	Error        error
	URLID        string
	Body         []byte
	ResponseTime time.Duration
	StatusCode   int
	ResponseSize int
}

func (s *Source) CheckHTTPURLs() ([]*URL, error) {
	conn, err := s.pool.Acquire(context.Background())
	if err != nil {
		return nil, errors.Wrap(err, "acquire connection")
	}
	defer conn.Release()

	ctx, cancel := context.WithTimeout(context.Background(), s.queryTimeout)
	defer cancel()

	rows, err := conn.Query(ctx, "select id, url from host_url")
	if err != nil {
		return nil, errors.Wrap(err, "execute")
	}

	var result []*URL
	for rows.Next() {
		u := URL{}
		err := rows.Scan(&u.ID, &u.URL)
		if err != nil {
			return nil, errors.Wrap(err, "scan result")
		}

		result = append(result, &u)
	}

	return result, nil
}

func (s *Source) CheckHTTPSaveResponse(r *HTTPResponse) error {
	conn, err := s.pool.Acquire(context.Background())
	if err != nil {
		return errors.Wrap(err, "acquire connection")
	}
	defer conn.Release()

	// 1. save to log and get the z-score
	// 2. alerting: save to alert log if z-score >= 2 and is_alerting == false and response_time > threshold
	// 3. unalerting: save to alert log if z-score < 2 and is_alerting == true and response_time <= threshold

	ctx, cancel := context.WithTimeout(context.Background(), s.queryTimeout)
	defer cancel()

	err = pgx.BeginFunc(ctx, conn, func(tx pgx.Tx) error {
		if _, err := tx.Exec(ctx, `
		insert into check_httprequest (
		  url_id,
		  agent_id,
		  response_time,
		  response_status_code,
		  response_size,
		  error) VALUES 
	($1, $2, $3, $4, $5, $6)
on CONFLICT (url_id, agent_id) do update set updated_at = now(), 
response_time = $3, response_status_code = $4, response_size = $5, error = $6
		`,
			r.URLID,
			s.agentID,
			r.ResponseTime/time.Millisecond,
			r.StatusCode,
			// r.Body,
			r.ResponseSize,
			r.Error,
		); err != nil {
			return err
		}

		// var zScore int
		if _, err := tx.Exec(ctx, `
		/*
		with z_score as (
			select 
				AVG(t.response_time) mean,
				stddev(t.response_time) sd
			from (
					select response_time 
					from check_httprequestlog 
					where url_id = $1 and agent_id = $2
					order by created_at limit 10
			) as t
		)
		*/
		insert into check_httprequestlog (url_id, agent_id, response_time, response_status_code, response_size, error) 
		select $1, $2, $3, $4, $5, $6
		-- from z_score
		RETURNING response_z_score
		`,
			r.URLID,
			s.agentID,
			r.ResponseTime/time.Millisecond,
			r.StatusCode,
			r.ResponseSize,
			r.Error,
		); err != nil {
			// fmt.Printf("`%v` %+v\n", err, r)
			return err
		}
		// else {
		// 	// var zScore int
		// 	// fmt.Println("row.String()", row.String())
		// 	// if zScore, err = strconv.Atoi(row.String()); err != nil {
		// 	// 	return err
		// 	// }
		// 	fmt.Println("z-score", zScore)
		// }

		return nil
	})

	if err != nil {
		return errors.Wrap(err, "tx")
	}

	return nil
}
