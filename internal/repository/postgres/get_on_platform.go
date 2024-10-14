package postgres

import (
	"context"
	"fmt"
	"session_manager/internal/domain"
	"session_manager/internal/domain/response"
	"time"
)

type UserActivity interface {
	GetUserActivityByMonth(ctx context.Context, dto *domain.UserActivity) (*response.UserActivity, error)
	GetUserActivityByDate(ctx context.Context, dto *domain.UserActivity) (*response.UserActivity, error)
}

func (s *storage) GetUserActivityByMonth(ctx context.Context, dto *domain.UserActivity) (*response.UserActivity, error) {
	ctx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()

	rows, err := s.pool.Query(ctx,
		`WITH monthly_hours AS (
			SELECT
				login,
				EXTRACT(YEAR FROM start_date_time) AS year,
				EXTRACT(MONTH FROM start_date_time) AS month_number,
				EXTRACT(EPOCH FROM (end_date_time - start_date_time)) / 3600 AS hours_calc
			FROM session.on_platform
			WHERE
				login = $1
				AND session_type = $2
				AND DATE_TRUNC('day', start_date_time) >= $3::date
				AND DATE_TRUNC('day', end_date_time) <= $4::date
		)
		SELECT 
			login,
			year,
			month_number,
			SUM(hours_calc) AS total_hours,
			SUM(SUM(hours_calc)) OVER (PARTITION BY login) AS total_hours
		FROM monthly_hours
		GROUP BY login, year, month_number
		ORDER BY year DESC, month_number DESC;`,
		dto.Login,
		dto.SessionType,
		dto.FromDate,
		dto.ToDate,
	)

	if err != nil {
		return nil, fmt.Errorf("query: %w", err)
	}

	defer rows.Close()

	return iterateRowsActivityByMonth(rows)
}

func (s *storage) GetUserActivityByDate(ctx context.Context, dto *domain.UserActivity) (*response.UserActivity, error) {
	ctx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()

	rows, err := s.pool.Query(ctx,
		`SELECT
			login,
			DATE_TRUNC('day', start_date_time) AS date,
			SUM(EXTRACT(EPOCH FROM (end_date_time - start_date_time)) / 3600) AS hours,
			SUM(EXTRACT(EPOCH FROM (end_date_time - start_date_time)) / 3600) OVER (PARTITION BY login) AS total_hours
		FROM session.on_platform
		WHERE
			login = $1
			AND session_type = $2
			AND DATE_TRUNC('day', start_date_time) >= $3::date
			AND DATE_TRUNC('day', end_date_time) <= $4::date
		GROUP BY login, date, start_date_time, end_date_time
		ORDER BY date;`,
		dto.Login,
		dto.SessionType,
		dto.FromDate,
		dto.ToDate,
	)

	if err != nil {
		return nil, fmt.Errorf("query: %w", err)
	}

	defer rows.Close()

	return iterateRowsActivityByDate(rows)
}
