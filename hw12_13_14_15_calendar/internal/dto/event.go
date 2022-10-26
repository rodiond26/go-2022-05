package dto

type EventDto struct {
	ID               int64  `json:"id"`
	Title            string `json:"title"`
	StartDate        string `json:"startDate"`
	EndDate          string `json:"endDate"`
	Description      string `json:"description"`
	UserID           int64  `json:"userId"`
	NotificationDate string `json:"notificationDate"`
}

type CreateEventRequest struct {
	Event *EventDto `json:"event"`
}

type CreateEventResponse struct {
	ID string `json:"id"`
}

type EventUpdateRequest struct {
	Event *EventDto `json:"event"`
}

type GetAllEventsResponse struct {
	Events []*EventDto `json:"events"`
}

type FindEventByIDResponse struct {
	Event *EventDto `json:"event"`
}

type FindEventsByDayResponse struct {
	Events []*EventDto `json:"events"`
}

type FindEventsByWeekResponse struct {
	Events []*EventDto `json:"events"`
}

type FindEventsByMonthResponse struct {
	Events []*EventDto `json:"events"`
}
