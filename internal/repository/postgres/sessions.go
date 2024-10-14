package postgres

import (
	"context"
	"fmt"
	"log"
	"session_manager/internal/domain"
	"session_manager/internal/domain/response"
	"time"

	"github.com/jackc/pgx/v5"
)

type Sessions interface {
	CreateSessionOnCampus(ctx context.Context, dto *domain.Campus) error
	CreateSessionOnPlatform(ctx context.Context, dto *domain.Platform) error
	IsSessionExists(ctx context.Context, login string) ([]response.Session, error)
}

func (s *storage) CreateSessionOnCampus(ctx context.Context, dto *domain.Campus) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	// start session
	if _, err := s.pool.Exec(ctx,
		`INSERT INTO 
		session.on_campus (id, comp_name, ip_addr, login, next_ping_sec, start_date_time, end_date_time) 
		VALUES ($1, $2, $3, $4, $5, $6, $7);`,
		dto.ID,
		dto.ComputerName,
		dto.IPAddress,
		dto.Login,
		int(dto.NextPing.Seconds()),
		dto.StartDateTime,
		dto.EndDateTime,
	); err != nil {
		return s.customErr("exec", err, dto)
	}

	return nil
}

func (s *storage) CreateSessionOnPlatform(ctx context.Context, dto *domain.Platform) error {
	ctx2, cancel := context.WithTimeout(ctx, 20*time.Second)
	defer cancel()

	updateSessionEndQuery := `UPDATE session.on_campus
	SET end_date_time = $1
	WHERE id = $2;`

	// -------------- if only on_campus
	if dto.SessionType == "" {
		if tag, err := s.pool.Exec(ctx2, updateSessionEndQuery,
			dto.EndDateTime,
			dto.SessionID,
		); err != nil {
			return s.customErr("on_campus: exec: update", err, dto)
		} else if tag.RowsAffected() == 0 {
			return response.ErrNotFound
		}
		return nil
	}

	// -------------- if other activity [on zero platforn and etc...]
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer func() {
		if err = tx.Rollback(ctx); err != nil {
			log.Printf("CreateSessionOnPlatform: rollback: %s", err.Error())
		}
	}()

	// start activity
	if _, err := tx.Exec(ctx2,
		`INSERT INTO session.on_platform (session_id, session_type, login, start_date_time, end_date_time)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (session_id, session_type)
		DO UPDATE SET
		end_date_time = EXCLUDED.end_date_time;`,
		dto.SessionID,
		dto.SessionType,
		dto.Login,
		dto.StartDateTime,
		dto.EndDateTime,
	); err != nil {
		return s.customErr("on_platform: exec: insert", err, dto)
	}

	// also update session end_date_time
	if tag, err := tx.Exec(ctx2, updateSessionEndQuery,
		dto.EndDateTime,
		dto.SessionID,
	); err != nil {
		return s.customErr("on_platform: exec: update", err, dto)
	} else if tag.RowsAffected() == 0 {
		return response.ErrNotFound
	}

	err = tx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("commit tx: %w", err)
	}

	return nil
}

func (s *storage) IsSessionExists(ctx context.Context, login string) ([]response.Session, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	rows, err := s.pool.Query(ctx,
		`SELECT c.id, c.comp_name, c.ip_addr, c.login, u.multi, c.start_date_time, c.end_date_time
		FROM session.on_campus c
		INNER JOIN env_tracker.users u
		ON c.login = u.login
		WHERE c.login = $1
		AND c.end_date_time >= NOW();`,
		login,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
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
			&session.Multi,
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
