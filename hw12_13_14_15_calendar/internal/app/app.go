package app

import (
	"context"
	"time"

	storage "github.com/rodiond26/go-2022-05/hw12_13_14_15_calendar/internal/storage"
	init_storage "github.com/rodiond26/go-2022-05/hw12_13_14_15_calendar/internal/storage/initializing"
)

type App struct {
	logger  Logger
	storage storage.Storage
}

type Event struct {
	ID               int64     `json:"id"`
	Title            string    `json:"title"`
	StartDate        time.Time `json:"startDate"`
	EndDate          time.Time `json:"endDate"`
	Description      string    `json:"description"`
	UserID           int64     `json:"userId"`
	NotificationDate time.Time `json:"notificationDate"`
}

type Logger interface {
	Info(msg string)
	Warn(msg string)
	Error(msg string)
	Debug(msg string)
}

type Storage interface { // TODO
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

func (a *App) CreateEvent(ctx context.Context, event *Event) (id int64, err error) {
	return a.storage.CreateEvent(ctx, &storage.Event{
		ID:               event.ID,
		Title:            event.Title,
		StartDate:        event.StartDate,
		EndDate:          event.EndDate,
		Description:      event.Description,
		UserID:           event.UserID,
		NotificationDate: event.NotificationDate,
	})
}
