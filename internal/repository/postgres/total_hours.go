package postgres

import (
	"context"
	"fmt"
	"session_manager/internal/domain"

	"time"

	"github.com/jackc/pgx/v5"
)

type CadetTotalHours interface {
	GetTotalHours(ctx context.Context, req *domain.CadetTotalHoursRequest) ([]domain.CadetTotalHoursResponse, error)
}

func (s *storage) GetTotalHours(ctx context.Context, req *domain.CadetTotalHoursRequest) ([]domain.CadetTotalHoursResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()

	rows, err := s.pool.Query(ctx,
		`WITH monthly_hours AS (
			SELECT
				u.id_zero,
				c.login,
				EXTRACT(YEAR FROM c.start_date_time) AS year,
				EXTRACT(MONTH FROM c.start_date_time) AS month_number,
				EXTRACT(EPOCH FROM (c.end_date_time - c.start_date_time)) / 3600 AS hours_calc
		FROM session.on_campus c
		JOIN env_tracker.users u ON c.login = u.login 
		WHERE
			u.div_id = $1
			AND DATE_TRUNC('day', c.start_date_time) >= $2::date
			AND DATE_TRUNC('day', c.end_date_time) <= $3::date
		)
		SELECT 
			id_zero,
			login,
			year,
			month_number,
			SUM(hours_calc) AS hours,
			SUM(SUM(hours_calc)) OVER (PARTITION BY login) AS total_hours
		FROM monthly_hours
		GROUP BY id_zero, login, year, month_number
		ORDER BY year DESC, month_number DESC;`,
		req.ModuleID,
		req.FromDate,
		req.ToDate,
	)

	if err != nil {
		return nil, fmt.Errorf("query: %w", err)
	}

	defer rows.Close()

	return iterateRowsTotalHours(rows)
}

func iterateRowsTotalHours(rows pgx.Rows) ([]domain.CadetTotalHoursResponse, error) {
	th := make([]domain.CadetTotalHoursResponse, 0, 1000)

	for rows.Next() {
		month := domain.CadetTotalHoursResponse{}
		if err := rows.Scan(
			&month.ID,
			&month.Login,
			&month.Year,
			&month.MonthNumber,
			&month.Hours,
			&month.TotalHours,
		); err != nil {
			return nil, fmt.Errorf("in iterate row: %w", err)
		}
		th = append(th, month)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows at all: %w", err)
	}

	return th, nil
}
