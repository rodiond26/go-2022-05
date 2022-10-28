package sqlstorage

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/rodiond26/go-2022-05/hw12_13_14_15_calendar/internal/model"
)

var (
	ErrIsCanceled       = errors.New("search is canceled")
	ErrIsBusy           = errors.New("time is busy")
	ErrEventIsNotFound  = errors.New("event is not found")
	years, months, days int
)

const (
	invalidID = -1
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

	config.MaxConns = int32(10)
	config.MinConns = int32(1)
	config.HealthCheckPeriod = 1 * time.Minute
	config.MaxConnLifetime = 1 * time.Hour
	config.MaxConnIdleTime = 30 * time.Minute
	config.ConnConfig.ConnectTimeout = 1 * time.Second
	pool, err = pgxpool.ConnectConfig(ctx, config)
	if err != nil {
		err = fmt.Errorf("when connecting using config [%+v] then error: [%w]", config, err)
		return
	}
	return
}

func (s *Storage) AddEvent(ctx context.Context, newEvent *model.Event) (id int64, err error) {
	events, err := s.FindEventsByPeriod(ctx, newEvent.StartDate, newEvent.EndDate)
	if err != nil {
		return invalidID, err
	}
	for _, event := range events {
		if newEvent.StartDate.After(event.StartDate) && newEvent.StartDate.Before(event.EndDate) {
			return invalidID, ErrIsBusy
		}
		if newEvent.EndDate.After(event.StartDate) && newEvent.EndDate.Before(event.EndDate) {
			return invalidID, ErrIsBusy
		}
	}

	query := `INSERT INTO events(event_id, title, start_date, end_date, description, user_id, remind_date)
	          VALUES ($1, $2, $3, $4, $5, $6, $7)
			  RETURNING event_id;`
	_, err = s.PgxPool.Exec(ctx, query, newEvent.ID, newEvent.Title, newEvent.StartDate,
		newEvent.EndDate, newEvent.Description, newEvent.UserID, newEvent.NotificationDate)
	if err != nil {
		return invalidID, err
	}
	return newEvent.ID, nil
}

func (s *Storage) FindEventByID(ctx context.Context, id int64) (event model.Event, err error) {
	query := `SELECT event_id, title, start_date, end_date, description, user_id, remind_date
	            FROM events
               WHERE event_id = $1;`

	var events []*model.Event
	err = pgxscan.Select(ctx, s.PgxPool, &events, query, id)
	if err != nil {
		return event, err
	}
	return *events[0], nil
}

func (s *Storage) UpdateEvent(ctx context.Context, event *model.Event) (err error) {
	query := `UPDATE events
	             SET title=$1, start_date=$2, end_date=$3, description=$4, user_id=$5, remind_date=$6
	           WHERE id=$7;`
	_, err = s.PgxPool.Exec(ctx, query, event.ID, event.Title, event.StartDate,
		event.EndDate, event.Description, event.UserID, event.NotificationDate, event.ID)
	if err != nil {
		return err
	}
	return nil
}

func (s *Storage) DeleteEventByID(ctx context.Context, id int64) (err error) {
	query := `DELETE FROM events 
	           WHERE id=$1`
	_, err = s.PgxPool.Exec(ctx, query, id)
	if err != nil {
		return err
	}
	return nil
}

func (s *Storage) FindEventsByPeriod(ctx context.Context, start, end time.Time) (events []model.Event, err error) {
	query := `SELECT event_id, title, start_date, end_date, description, user_id, remind_date
		        FROM events
		       WHERE start_date BETWEEN $1 AND $2;`
	rows, err := s.PgxPool.Query(ctx, query, start, end)
	if err != nil {
		err = fmt.Errorf("when executing query [%s] then error: [%w]", query, err)
		return events, err
	}
	defer rows.Close()

	for rows.Next() {
		values, rowErr := rows.Values()
		if rowErr != nil {
			err = fmt.Errorf("when iterating dataset then error: [%w]", rowErr)
			return events, err
		}
		event := model.Event{
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

func (s *Storage) Close(ctx context.Context) (err error) {
	s.PgxPool.Close()
	return nil
}

func (s *Storage) FindEventsByDay(ctx context.Context, start time.Time) (events []model.Event, err error) {
	years = 0
	months = 0
	days = 1
	end := start.AddDate(years, months, days)

	return s.FindEventsByPeriod(ctx, start, end)
}

func (s *Storage) FindEventsByWeek(ctx context.Context, start time.Time) (events []model.Event, err error) {
	years = 0
	months = 0
	days = 7
	end := start.AddDate(years, months, days)

	return s.FindEventsByPeriod(ctx, start, end)
}

func (s *Storage) FindEventsByMonth(ctx context.Context, start time.Time) (events []model.Event, err error) {
	years = 0
	months = 1
	days = 0
	end := start.AddDate(years, months, days)

	return s.FindEventsByPeriod(ctx, start, end)
}
