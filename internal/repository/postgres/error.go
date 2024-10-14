package postgres

import (
	"context"
	"fmt"
	"log"
	"session_manager/internal/domain"
	"session_manager/internal/domain/response"
	"time"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
)

func (s *storage) customErr(message string, err error, jsn domain.DataToJson) error {
	// postgres errors
	if pgErr, ok := err.(*pgconn.PgError); ok {
		go s.writeErr(fmt.Errorf("%s: %s", message, err), jsn)
		// duplicate key
		if pgErr.Code == pgerrcode.UniqueViolation {
			return response.ErrDuplicateKey
		}
		// unknown keys
		if pgErr.Code == pgerrcode.ForeignKeyViolation {
			return response.ErrForeignKey
		}
		// custom errors from trigger
		if pgErr.Code == pgerrcode.RaiseException {
			if pgErr.Message == response.ErrEndStartDate.Error() {
				return response.ErrEndStartDate
			}
			if pgErr.Message == response.ErrEndEndDate.Error() {
				return response.ErrEndEndDate
			}
		}
	}
	// service errors
	if message != "" {
		return fmt.Errorf("%s: %s", message, err)
	}
	return err
}

func (s *storage) writeErr(err error, jsn domain.DataToJson) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	// start activity
	if _, err := s.pool.Exec(ctx,
		`INSERT INTO session.errors (unique_key, error_body, request_data, created)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (unique_key)
		DO UPDATE SET
		error_body = EXCLUDED.error_body,
		request_data = EXCLUDED.request_data,
		created = EXCLUDED.created;`,
		jsn.Key(),
		err.Error(),
		jsn.Marshal(),
		time.Now(),
	); err != nil {
		log.Printf("write error to db: exec: insert %s", err)
	}

	log.Println("the error has been logged into the database")
}
