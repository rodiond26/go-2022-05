package app

import (
	"context"
	"fmt"
	"time"

	"github.com/rodiond26/go-2022-05/hw12_13_14_15_calendar/internal/model"
	"github.com/rodiond26/go-2022-05/hw12_13_14_15_calendar/internal/storage"
	init_storage "github.com/rodiond26/go-2022-05/hw12_13_14_15_calendar/internal/storage/initializing"
)

type App struct {
	logger  Logger
	storage storage.Storage
}

type Logger interface {
	Info(msg string)
	Warn(msg string)
	Error(msg string)
	Debug(msg string)
}

type Application interface {
	Close(ctx context.Context) error
	AddEvent(ctx context.Context, newEvent *model.Event) (id int64, err error)
	FindEventByID(ctx context.Context, id int64) (event model.Event, err error)
	UpdateEvent(ctx context.Context, event *model.Event) (err error)
	DeleteEventByID(ctx context.Context, id int64) (err error)
	FindEventsByPeriod(ctx context.Context, startDate, endDate time.Time) (events []model.Event, err error)
	FindEventsByDay(ctx context.Context, startDate time.Time) (events []model.Event, err error)
	FindEventsByWeek(ctx context.Context, startDate time.Time) (events []model.Event, err error)
	FindEventsByMonth(ctx context.Context, startDate time.Time) (events []model.Event, err error)
}

func New(logger Logger, storage storage.Storage) *App {
	return &App{
		logger:  logger,
		storage: storage,
	}
}

func (a *App) Connect(ctx context.Context, storageType string, dsn string) error {
	storage, err := init_storage.NewStorage(ctx, storageType, dsn)
	if err != nil {
		return err
	}
	a.storage = storage
	return nil
}

func (a *App) Close(ctx context.Context) error {
	return a.storage.Close(ctx)
}

func (a *App) AddEvent(ctx context.Context, newEvent *model.Event) (id int64, err error) {
	return a.storage.AddEvent(ctx, newEvent)
}

func (a *App) UpdateEvent(ctx context.Context, event *model.Event) (err error) {
	if event.Title == "" {
		return fmt.Errorf("empty title of event")
	}
	if event.EndDate.Before(event.StartDate) {
		return fmt.Errorf("wrong dates")
	}
	return a.storage.UpdateEvent(ctx, event)
}

func (a *App) FindEventByID(ctx context.Context, id int64) (event model.Event, err error) {
	event, err = a.storage.FindEventByID(ctx, id)
	if err != nil {
		return event, err
	}
	return event, nil
}

func (a *App) DeleteEventByID(ctx context.Context, id int64) (err error) {
	return a.storage.DeleteEventByID(ctx, id)
}

func (a *App) FindEventsByPeriod(ctx context.Context, startDate time.Time) ([]model.Event, error) {
	return a.storage.FindEventsByPeriod(ctx, startDate, startDate)
}

func (a *App) FindEventsByDay(ctx context.Context, startDate time.Time) (events []model.Event, err error) {
	return a.storage.FindEventsByDay(ctx, startDate)
}

func (a *App) FindEventsByWeek(ctx context.Context, startDate time.Time) (events []model.Event, err error) {
	return a.storage.FindEventsByWeek(ctx, startDate)
}

func (a *App) FindEventsByMonth(ctx context.Context, startDate time.Time) (events []model.Event, err error) {
	return a.storage.FindEventsByMonth(ctx, startDate)
}
