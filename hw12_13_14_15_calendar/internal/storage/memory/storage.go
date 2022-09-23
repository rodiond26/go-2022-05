package memorystorage

import (
	"context"
	"sync"
	"time"

	"github.com/rodiond26/go-2022-05/hw12_13_14_15_calendar/internal/storage"
)

type Storage struct {
	mu     sync.RWMutex
	events map[int64]storage.Event
}

func New() *Storage {
	return &Storage{
		mu:     sync.RWMutex{},
		events: make(map[int64]storage.Event),
	}
}

func (s *Storage) Close(ctx context.Context) error {
	return nil
}

func (s *Storage) FindEventsByPeriod(ctx context.Context, startDate, endDate time.Time) (events []storage.Event, err error) {

	select {
	case <-ctx.Done():
		return nil, nil // TODO fix it
	default:
		s.mu.Lock()
		defer s.mu.Unlock()
// TODO fix it
		for _, event := range s.events {
			if startDate.Before(event.StartDate) && endDate.After(event.StartDate) {
				events = append(events, event)
			}
		}
	}
	return
}
