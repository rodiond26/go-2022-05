package storage

import (
	"context"
	"time"
)

type Event struct {
	ID               int64
	Title            string
	StartDate        time.Time
	EndDate          time.Time
	Description      string
	UserID           int64
	NotificationDate time.Time
}

type User struct {
	ID          int64
	Name        string
	Description string
}

type Storage interface {
	CreateEvent(ctx context.Context, newEvent *Event) (id int64, err error)
	FindEventByID(ctx context.Context, id int64) (event Event, err error)
	UpdateEvent(ctx context.Context, event *Event) (err error)
	DeleteEventByID(ctx context.Context, id int64) (err error)

	FindEventsByPeriod(ctx context.Context, startDate, endDate time.Time) (events []Event, err error)
	Close(ctx context.Context) (err error)
}
