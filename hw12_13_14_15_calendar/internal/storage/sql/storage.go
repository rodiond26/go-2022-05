package sqlstorage

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/pkg/errors"
)

type Storage struct {
	db *sql.DB
}

func New() *Storage {
	return &Storage{
		db: nil,
	}
}

func (s *Storage) Connect(ctx context.Context, dsn string) (err error) {
	s.db, err = sql.Open("pgx", dsn)
	if err != nil {
		return fmt.Errorf("cannot open pgx driver: %w", err)
	}

	return s.db.PingContext(ctx)
}

func (s *Storage) Close(ctx context.Context) error {
	return s.db.Close()
}

func (s *Storage) GetEventsByPeriod(startDate time.Time, endDate time.Time) ([]Event, error) {
	events, err := s.db.Query(
		`SELECT event_id,
       			title,
       			start_date,
    		    end_date,
    		    description,
    		    user_id,
    		    remind_in
		   FROM event
		  WHERE start_date >=$1
		    AND start_date <=$2`,
		startDate,
		endDate,
	)
	if err != nil {
		return nil, errors.Wrap(err, "cannot execute query")
	}
	defer events.Close()

	return events, nil
}
