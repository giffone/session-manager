package postgres

import (
	"github.com/jackc/pgx/v5/pgxpool"
)

type Storage interface {
	Sessions
	Dashboards
	UserActivity
	UserActivityInCampus
	CadetTotalHours
}

func NewStorage(pool *pgxpool.Pool) Storage {
	return &storage{pool: pool}
}

type storage struct {
	pool *pgxpool.Pool
}
