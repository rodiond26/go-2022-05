package model

import "time"

type Event struct {
	ID               int64     `json:"id"`
	Title            string    `json:"title"`
	StartDate        time.Time `json:"startDate"`
	EndDate          time.Time `json:"endDate"`
	Description      string    `json:"description"`
	UserID           int64     `json:"userId"`
	NotificationDate time.Time `json:"notificationDate"`
}
