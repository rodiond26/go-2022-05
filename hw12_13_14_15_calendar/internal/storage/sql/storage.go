package sqlstorage

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"

	"github.com/rodiond26/go-2022-05/hw12_13_14_15_calendar/internal/storage"
)

type Storage struct {
	PgxPool *pgxpool.Pool
}

func New() *Storage {
	return &Storage{
		PgxPool: nil,
	}
}

func (s *Storage) Connect(ctx context.Context, dsn string) (pool *pgxpool.Pool, err error) {
	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		err = fmt.Errorf("when parsing config [%s] then error: [%w]", dsn, err)
		return
	}

	// TODO add config parameters
	pool, err = pgxpool.ConnectConfig(ctx, config)
	if err != nil {
		err = fmt.Errorf("when connecting using config [%+v] then error: [%w]", config, err)
		return
	}
	return
}

func (s *Storage) Close(ctx context.Context) (err error) {
	s.PgxPool.Close()
	return nil
}

func (s *Storage) FindEventsByPeriod(ctx context.Context, startDate, endDate time.Time) (events []storage.Event, err error) {
	query := `SELECT event_id, title, start_date, end_date, description, user_id, remind_date
		        FROM event
		       WHERE start_date BETWEEN $1 AND $2`
	rows, err := s.PgxPool.Query(ctx, query, startDate, endDate)
	defer rows.Close()
	if err != nil {
		err = fmt.Errorf("when executing query [%s] then error: [%w]", query, err)
		return events, err
	}

	for rows.Next() {
		values, rowErr := rows.Values()
		if rowErr != nil {
			err = fmt.Errorf("when iterating dataset then error: [%w]", rowErr)
			return events, err
		}
		event := storage.Event{
			ID:               values[0].(int64),
			Title:            values[1].(string),
			StartDate:        values[2].(time.Time),
			EndDate:          values[3].(time.Time),
			Description:      values[4].(string),
			UserID:           values[5].(int64),
			NotificationDate: values[6].(time.Time),
		}
		events = append(events, event)
	}

	return events, err
}
