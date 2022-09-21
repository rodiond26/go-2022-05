package storage

import "time"

type Event struct {
	ID               int64
	Title            string
	StartDate        time.Time
	EndDate          time.Time
	Description      string
	UserID           int64
	NotificationDate time.Time
}
