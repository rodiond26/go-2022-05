package storage

import (
	"context"
	"time"

	"github.com/rodiond26/go-2022-05/hw12_13_14_15_calendar/internal/model"
)

type User struct {
	ID          int64
	Name        string
	Description string
}

type Storage interface {
	AddEvent(ctx context.Context, newEvent *model.Event) (id int64, err error)
	FindEventByID(ctx context.Context, id int64) (event model.Event, err error)
	UpdateEvent(ctx context.Context, event *model.Event) (err error)
	DeleteEventByID(ctx context.Context, id int64) (err error)

	FindEventsByPeriod(ctx context.Context, startDate, endDate time.Time) (events []model.Event, err error)
	Close(ctx context.Context) (err error)
}
