package httpServer

import (
	"context"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/rodiond26/go-2022-05/hw12_13_14_15_calendar/internal/app"
)

type Server struct {
	httpServer *http.Server
	logger     Logger
	router     *mux.Router
	app        *app.App
}

type Logger interface {
	Info(msg string)
	Error(msg string)
	Warn(msg string)
	Debug(msg string)
}

func NewServer(logger Logger, app *app.App) *Server {
	s := &Server{
		logger: logger,
		router: mux.NewRouter(),
		app:    app,
	}
	s.router.Use(LoggingMiddleware(s.logger))

	s.router.HandleFunc("/events/day", s.findEventsByDay).Methods("GET")
	s.router.HandleFunc("/events/week", s.findEventsByWeek).Methods("GET")
	s.router.HandleFunc("/events/month", s.findEventsByMonth).Methods("GET")
	s.router.HandleFunc("/events/{id}", s.findEventByID).Methods("GET")
	s.router.HandleFunc("/events", s.addEvent).Methods("POST")
	s.router.HandleFunc("/events/{id}", s.updateEvent).Methods("PUT")
	s.router.HandleFunc("/events/{id}", s.deleteEventByID).Methods("DELETE")

	return s
}

func (s *Server) Start(ctx context.Context, addr string) error {
	s.logger.Info("HTTP server [" + addr + "] starting ...")
	s.httpServer = &http.Server{
		Addr:              addr,
		Handler:           s.router,
		WriteTimeout:      5 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       5 * time.Second,
	}

	errChan := make(chan error)

	go func() {
		if err := s.httpServer.ListenAndServe(); err != nil {
			errChan <- err
		}
	}()

	select {
	case <-ctx.Done():
	case err := <-errChan:
		return err
	}
	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	s.logger.Info("HTTP server is stopped ...")
	if err := s.httpServer.Shutdown(ctx); err != nil {
		return err
	}
	return nil
}
