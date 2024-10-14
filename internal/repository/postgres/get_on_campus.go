package postgres

import (
	"context"
	"fmt"
	"math"
	"session_manager/internal/domain"
	"session_manager/internal/domain/response"

	"time"

	"github.com/jackc/pgx/v5"
)

type UserActivityInCampus interface {
	GetUserActivityByMonthInCampus(ctx context.Context, dto *domain.UserActivity) (*response.UserActivity, error)
	GetUserActivityByDateInCampus(ctx context.Context, dto *domain.UserActivity) (*response.UserActivity, error)
}

func (s *storage) GetUserActivityByMonthInCampus(ctx context.Context, dto *domain.UserActivity) (*response.UserActivity, error) {
	ctx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()

	rows, err := s.pool.Query(ctx,
		`WITH monthly_hours AS (
			SELECT
				login,
				EXTRACT(YEAR FROM start_date_time) AS year,
				EXTRACT(MONTH FROM start_date_time) AS month_number,
				EXTRACT(EPOCH FROM (end_date_time - start_date_time)) / 3600 AS hours_calc
			FROM session.on_campus
			WHERE
				login = $1
				AND DATE_TRUNC('day', start_date_time) >= $2::date
				AND DATE_TRUNC('day', end_date_time) <= $3::date
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
		dto.FromDate,
		dto.ToDate,
	)

	if err != nil {
		return nil, fmt.Errorf("query: %w", err)
	}

	defer rows.Close()

	return iterateRowsActivityByMonth(rows)
}

func (s *storage) GetUserActivityByDateInCampus(ctx context.Context, dto *domain.UserActivity) (*response.UserActivity, error) {
	ctx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()

	rows, err := s.pool.Query(ctx,
		`SELECT
			login,
			DATE_TRUNC('day', start_date_time) AS date,
			SUM(EXTRACT(EPOCH FROM (end_date_time - start_date_time)) / 3600) AS hours,
			SUM(EXTRACT(EPOCH FROM (end_date_time - start_date_time)) / 3600) OVER (PARTITION BY login) AS total_hours
		FROM session.on_campus
		WHERE
			login = $1
			AND DATE_TRUNC('day', start_date_time) >= $2::date
			AND DATE_TRUNC('day', end_date_time) <= $3::date
		GROUP BY login, date, start_date_time, end_date_time
		ORDER BY date;`,
		dto.Login,
		dto.FromDate,
		dto.ToDate,
	)

	if err != nil {
		return nil, fmt.Errorf("query: %w", err)
	}

	defer rows.Close()

	return iterateRowsActivityByDate(rows)
}

func iterateRowsActivityByMonth(rows pgx.Rows) (*response.UserActivity, error) {
	activities := make([]response.UserActivityByMonth, 0, 36)
	var totalHours float64
	var login string

	for rows.Next() {
		activity := response.UserActivityByMonth{}
		if err := rows.Scan(
			&login,
			&activity.Year,
			&activity.MonthNumber,
			&activity.Hours,
			&totalHours,
		); err != nil {
			return nil, fmt.Errorf("in iterate row: %w", err)
		}
		activities = append(activities, activity)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows at all: %w", err)
	}

	return &response.UserActivity{
		Login:        login,
		TotalHours:   float32(math.Round(totalHours*100) / 100),
		UserActivity: activities,
	}, nil
}

func iterateRowsActivityByDate(rows pgx.Rows) (*response.UserActivity, error) {
	activities := make([]response.UserActivityByDate, 0, 360)
	var totalHours float64
	var login string

	for rows.Next() {
		activity := response.UserActivityByDate{}
		if err := rows.Scan(
			&login,
			&activity.Date,
			&activity.Hours,
			&totalHours,
		); err != nil {
			return nil, fmt.Errorf("in iterate row: %w", err)
		}
		activities = append(activities, activity)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows at all: %w", err)
	}

	return &response.UserActivity{
		Login:        login,
		TotalHours:   float32(math.Round(totalHours*100) / 100),
		UserActivity: activities,
	}, nil
}
