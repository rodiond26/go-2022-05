package memorystorage

import (
	"context"
	"errors"
	"sync"
	"time"

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
	mu     sync.RWMutex
	events map[int64]model.Event
}

func New() *Storage {
	return &Storage{
		mu:     sync.RWMutex{},
		events: make(map[int64]model.Event),
	}
}

func (s *Storage) AddEvent(ctx context.Context, newEvent *model.Event) (id int64, err error) {
	select {
	case <-ctx.Done():
		return invalidID, ErrIsCanceled
	default:
		s.mu.Lock()
		defer s.mu.Unlock()

		for _, event := range s.events {
			if newEvent.StartDate.After(event.StartDate) && newEvent.StartDate.Before(event.EndDate) {
				return invalidID, ErrIsBusy
			}
			if newEvent.EndDate.After(event.StartDate) && newEvent.EndDate.Before(event.EndDate) {
				return invalidID, ErrIsBusy
			}
		}
		id = newEvent.ID
		s.events[id] = *newEvent
	}
	return id, nil
}

func (s *Storage) FindEventByID(ctx context.Context, id int64) (event model.Event, err error) {
	select {
	case <-ctx.Done():
		return event, ErrIsCanceled
	default:
		s.mu.Lock()
		defer s.mu.Unlock()

		if _, ok := s.events[id]; ok {
			return s.events[id], nil
		}
	}
	return event, ErrEventIsNotFound
}

func (s *Storage) UpdateEvent(ctx context.Context, updatedEvent *model.Event) (err error) {
	select {
	case <-ctx.Done():
		return ErrIsCanceled
	default:
		s.mu.Lock()
		defer s.mu.Unlock()

		if _, ok := s.events[updatedEvent.ID]; !ok {
			return ErrEventIsNotFound
		}
		s.events[updatedEvent.ID] = *updatedEvent
	}
	return nil
}

func (s *Storage) DeleteEventByID(ctx context.Context, id int64) (err error) {
	select {
	case <-ctx.Done():
		return ErrIsCanceled
	default:
		s.mu.Lock()
		defer s.mu.Unlock()

		if _, ok := s.events[id]; !ok {
			return ErrEventIsNotFound
		}
		delete(s.events, id)
	}
	return nil
}

func (s *Storage) FindEventsByPeriod(ctx context.Context, start, end time.Time) (events []model.Event, err error) {
	events = make([]model.Event, 0)
	select {
	case <-ctx.Done():
		return nil, ErrIsCanceled
	default:
		s.mu.Lock()
		defer s.mu.Unlock()

		for _, event := range s.events {
			if start.Before(event.StartDate) && end.After(event.StartDate) {
				events = append(events, event)
			}
		}
	}
	return events, nil
}

func (s *Storage) Close(ctx context.Context) error {
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
