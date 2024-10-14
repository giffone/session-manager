package postgres

import (
	"context"
	"fmt"
	"session_manager/internal/domain/response"
	"time"
)

type Dashboards interface {
	GetOnlineDashboard(ctx context.Context) ([]response.Session, error)
}

func (s *storage) GetOnlineDashboard(ctx context.Context) ([]response.Session, error) {
	ctx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()

	rows, err := s.pool.Query(ctx,
		`SELECT id, comp_name, ip_addr, login, start_date_time, end_date_time
		FROM session.on_campus
		WHERE end_date_time >= (NOW() - INTERVAL '10 seconds');`,
	)
	if err != nil {
		return nil, fmt.Errorf("query: %w", err)
	}
	defer rows.Close()

	sessions := make([]response.Session, 0, 250)

	for rows.Next() {
		session := response.Session{}
		if err := rows.Scan(
			&session.ID,
			&session.ComputerName,
			&session.IPAddress,
			&session.Login,
			&session.StartDateTime,
			&session.EndDateTime,
		); err != nil {
			return nil, fmt.Errorf("in iterate row: %w", err)
		}
		sessions = append(sessions, session)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows at all: %w", err)
	}

	return sessions, nil
}
