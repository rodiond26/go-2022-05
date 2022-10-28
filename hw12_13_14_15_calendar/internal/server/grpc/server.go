//go:generate protoc -I ../../../api EventService.proto --go_out=. --go-grpc_out=.
package grpc

import (
	context "context"
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/rodiond26/go-2022-05/hw12_13_14_15_calendar/internal/app"
	"github.com/rodiond26/go-2022-05/hw12_13_14_15_calendar/internal/logger"
	"github.com/rodiond26/go-2022-05/hw12_13_14_15_calendar/internal/model"
	"google.golang.org/grpc"
)

const timeLayout = "2006.01.02 15:04:05"

type Server struct {
	grpcServer *grpc.Server
	logger     *logger.Logger
	app        *app.App
}

func NewServer(logger *logger.Logger, app *app.App) *Server {
	return &Server{
		app:    app,
		logger: logger,
	}
}

func (s *Server) Start(ctx context.Context, addr string) error {
	s.logger.Info("GRPC server is starting " + addr)
	lsn, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	s.grpcServer = grpc.NewServer(grpc.UnaryInterceptor(loggingServerInterceptor(*s.logger)))
	RegisterCalendarServer(s.grpcServer, s)
	if err = s.grpcServer.Serve(lsn); err != nil {
		return err
	}
	return nil
}

func (s *Server) mustEmbedUnimplementedCalendarServer() {
}

func (s *Server) Stop(ctx context.Context) error {
	s.logger.Info("GRPC server is stopping...")
	s.grpcServer.GracefulStop()
	return nil
}

func loggingServerInterceptor(logger app.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (_ interface{}, err error) { //nolint:lll
		logger.Info(fmt.Sprintf("method: %s, duration: %s, request: %+v", info.FullMethod, time.Since(time.Now()), req))
		i, err := handler(ctx, req)
		return i, err
	}
}

func (s *Server) GetEventByID(ctx context.Context, e *GetEventByIDRequest) (*GetEventByIDResponse, error) {
	id := e.Id
	event, err := s.app.FindEventByID(ctx, id)
	if err != nil {
		return nil, err
	}
	res := marshalEvent(&event)
	if err != nil {
		return nil, err
	}
	return &GetEventByIDResponse{Event: res}, nil
}

func (s *Server) CreateEvent(ctx context.Context, e *AddEventRequest) (*AddEventResponse, error) {
	event, err := unmarshalEvent(e.Event)
	if err != nil {
		return nil, err
	}
	id, err := s.app.AddEvent(ctx, event)
	if err != nil {
		return nil, err
	}
	return &AddEventResponse{
		Id: id,
	}, nil
}

func (s *Server) UpdateEvent(ctx context.Context, e *EventUpdateRequest) (*empty.Empty, error) {
	event, err := unmarshalEvent(e.Event)
	if err != nil {
		return nil, err
	}
	err = s.app.UpdateEvent(ctx, event)
	if err != nil {
		return nil, err
	}
	return &empty.Empty{}, nil
}

func (s *Server) DeleteEvent(ctx context.Context, e *DeleteEventRequest) (*empty.Empty, error) {
	id := e.Id
	err := s.app.DeleteEventByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return &empty.Empty{}, nil
}

func (s *Server) FindEventsByDay(ctx context.Context, dayDate *FindDayEventsRequest) (*FindDayEventsResponse, error) {
	d := dayDate.StartDate
	start, err := unmarshalStringToTime(d)
	if err != nil {
		return nil, err
	}
	events, err := s.app.FindEventsByDay(ctx, start)
	if err != nil {
		return nil, err
	}
	res := marshalEvents(events)
	return &FindDayEventsResponse{Events: res}, nil
}

func (s *Server) FindEventsByWeek(ctx context.Context, weekDay *FindWeekEventsRequest) (*FindWeekEventsResponse, error) { //nolint:lll
	d := weekDay.StartDate
	start, err := unmarshalStringToTime(d)
	if err != nil {
		return nil, err
	}
	events, err := s.app.FindEventsByWeek(ctx, start)
	if err != nil {
		return nil, err
	}
	res := marshalEvents(events)
	return &FindWeekEventsResponse{Events: res}, nil
}

func (s *Server) FindEventsByMonth(ctx context.Context, monthDay *FindMonthEventsRequest) (*FindMonthEventsResponse, error) { //nolint:lll
	d := monthDay.StartDate
	start, err := unmarshalStringToTime(d)
	if err != nil {
		return nil, err
	}
	events, err := s.app.FindEventsByMonth(ctx, start)
	if err != nil {
		return nil, err
	}
	res := marshalEvents(events)
	return &FindMonthEventsResponse{Events: res}, nil
}

func unmarshalEvent(e *Event) (event *model.Event, err error) {
	startDate, err := unmarshalStringToTime(e.StartDate)
	if err != nil {
		return event, err
	}
	endDate, err := unmarshalStringToTime(e.EndDate)
	if err != nil {
		return event, err
	}
	notificationDate, err := unmarshalStringToTime(e.EndDate)
	if err != nil {
		return event, err
	}

	return &model.Event{
		ID:               e.Id,
		Title:            e.Title,
		StartDate:        startDate,
		EndDate:          endDate,
		Description:      e.Description,
		UserID:           e.UserId,
		NotificationDate: notificationDate,
	}, nil
}

func marshalEvent(e *model.Event) *Event {
	startDate := marshalTimeToString(e.StartDate)
	endDate := marshalTimeToString(e.EndDate)
	notificationDate := marshalTimeToString(e.EndDate)

	return &Event{
		Id:          e.ID,
		Title:       e.Title,
		StartDate:   startDate,
		EndDate:     endDate,
		Description: e.Description,
		UserId:      e.UserID,
		NotifyDate:  notificationDate,
	}
}

func unmarshalStringToTime(str string) (time.Time, error) {
	s := strings.Trim(str, `"`) // remove quotes
	return time.Parse(timeLayout, s)
}

func marshalTimeToString(t time.Time) (str string) {
	if t.IsZero() {
		return ""
	}
	return fmt.Sprintf(`"%s"`, t.Format(timeLayout))
}

func marshalEvents(e []model.Event) (events []*Event) {
	for _, event := range e {
		ev := event
		me := marshalEvent(&ev)
		events = append(events, me)
	}
	return events
}
