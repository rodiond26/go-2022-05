package httpServer

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/rodiond26/go-2022-05/hw12_13_14_15_calendar/internal/dto"
	"github.com/rodiond26/go-2022-05/hw12_13_14_15_calendar/internal/model"
)

const (
	ct         = "Content-Type"
	aj         = "application/json"
	timeLayout = "2006.01.02 15:04:05"
)

func (s *Server) findEventByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set(ct, aj)
	params := mux.Vars(r)
	eventID, err := strconv.ParseInt(params["id"], 10, 64)
	if err != nil {
		s.logger.Error(fmt.Sprintf("error to get request body: %v", err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	event, err := s.app.FindEventByID(r.Context(), eventID)
	if err != nil {
		s.logger.Error(fmt.Sprintf("error to get list of events: %v", err))
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	err = json.NewEncoder(w).Encode(event)
	if err != nil {
		s.logger.Error(fmt.Sprintf("error to get list of events: %v", err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (s *Server) addEvent(w http.ResponseWriter, r *http.Request) {
	w.Header().Set(ct, aj)
	eventReq := dto.CreateEventRequest{}
	err := json.NewDecoder(r.Body).Decode(&eventReq)
	if err != nil {
		s.logger.Error(fmt.Sprintf("failed to get body [%v]\n", err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ev, err := unmarshalEvent(eventReq.Event)
	if err != nil {
		s.logger.Error(fmt.Sprintf("failed to get body [%v]\n", err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	id, err := s.app.AddEvent(context.Background(), ev)
	if err != nil {
		s.logger.Error(fmt.Sprintf("failed to create event [%v]\n", err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = json.NewEncoder(w).Encode(&dto.CreateEventResponse{
		ID: fmt.Sprint(id),
	})

	if err != nil {
		s.logger.Error(fmt.Sprintf("error in sending response [%v]\n", err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (s *Server) updateEvent(w http.ResponseWriter, r *http.Request) {
	w.Header().Set(ct, aj)
	eventReq := dto.CreateEventRequest{}
	err := json.NewDecoder(r.Body).Decode(&eventReq)
	if err != nil {
		s.logger.Error(fmt.Sprintf("failed to get body [%v]\n", err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ev, err := unmarshalEvent(eventReq.Event)
	if err != nil {
		s.logger.Error(fmt.Sprintf("failed to get body [%v]\n", err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = s.app.UpdateEvent(r.Context(), ev)
	if err != nil {
		s.logger.Error(fmt.Sprintf("error to update event: %v", err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (s *Server) deleteEventByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set(ct, aj)
	params := mux.Vars(r)
	eventID, err := strconv.ParseInt(params["id"], 10, 64)
	if err != nil {
		s.logger.Error(fmt.Sprintf("error to get request body: %v", err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = s.app.DeleteEventByID(r.Context(), eventID)
	if err != nil {
		s.logger.Error(fmt.Sprintf("error to delete event: %v", err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (s *Server) findEventsByPeriod(w http.ResponseWriter, r *http.Request) {
	var day time.Time
	err := json.NewDecoder(r.Body).Decode(&day)
	if err != nil {
		s.logger.Error(fmt.Sprintf("error to get request body: %v", err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	events, err := s.app.FindEventsByPeriod(r.Context(), day)
	if err != nil {
		s.logger.Error(fmt.Sprintf("error to get list of events: %v", err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = json.NewEncoder(w).Encode(events)
	if err != nil {
		s.logger.Error(fmt.Sprintf("error to get list of events: %v", err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (s *Server) findEventsByDay(w http.ResponseWriter, r *http.Request) {
	timeNow := time.Now()

	events, err := s.app.FindEventsByDay(r.Context(), timeNow)
	if err != nil {
		s.logger.Error(fmt.Sprintf("error to get list of events: %v", err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = json.NewEncoder(w).Encode(events)
	if err != nil {
		s.logger.Error(fmt.Sprintf("error to get list of events: %v", err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (s *Server) findEventsByWeek(w http.ResponseWriter, r *http.Request) {
	timeNow := time.Now()

	events, err := s.app.FindEventsByWeek(r.Context(), timeNow)
	if err != nil {
		s.logger.Error(fmt.Sprintf("error to get list of events: %v", err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = json.NewEncoder(w).Encode(events)
	if err != nil {
		s.logger.Error(fmt.Sprintf("error to get list of events: %v", err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (s *Server) findEventsByMonth(w http.ResponseWriter, r *http.Request) {
	timeNow := time.Now()

	events, err := s.app.FindEventsByMonth(r.Context(), timeNow)
	if err != nil {
		s.logger.Error(fmt.Sprintf("error to get list of events: %v", err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = json.NewEncoder(w).Encode(events)
	if err != nil {
		s.logger.Error(fmt.Sprintf("error to get list of events: %v", err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func unmarshalEvent(e *dto.EventDto) (event *model.Event, err error) {
	startDate, err := unmarshalJSONToTime(e.StartDate)
	if err != nil {
		return event, err
	}
	endDate, err := unmarshalJSONToTime(e.EndDate)
	if err != nil {
		return event, err
	}
	notificationDate, err := unmarshalJSONToTime(e.EndDate)
	if err != nil {
		return event, err
	}

	return &model.Event{
		ID:               e.ID,
		Title:            e.Title,
		StartDate:        startDate,
		EndDate:          endDate,
		Description:      e.Description,
		UserID:           e.UserID,
		NotificationDate: notificationDate,
	}, nil
}

func marshalEvent(e *model.Event) *dto.EventDto {
	startDate := marshalTimeToJSON(e.StartDate)
	endDate := marshalTimeToJSON(e.EndDate)
	notificationDate := marshalTimeToJSON(e.EndDate)

	return &dto.EventDto{
		ID:               e.ID,
		Title:            e.Title,
		StartDate:        startDate,
		EndDate:          endDate,
		Description:      e.Description,
		UserID:           e.UserID,
		NotificationDate: notificationDate,
	}
}

func unmarshalJSONToTime(str string) (time.Time, error) {
	s := strings.Trim(str, `"`) // remove quotes
	return time.Parse(timeLayout, s)
}

func marshalTimeToJSON(t time.Time) (str string) {
	if t.IsZero() {
		return ""
	}
	return fmt.Sprintf(`"%s"`, t.Format(timeLayout))
}
