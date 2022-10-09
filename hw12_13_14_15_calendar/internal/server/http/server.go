package internalhttp

import (
	"context"
	"net/http"
	"time"
)

type Server struct {
	// Address string
	httpServer *http.Server
	// logger Logger
	// api     *app.App
}

type Logger interface {
	// 	Info(msg string)
	// 	Error(msg string)
	// 	Warn(msg string)
	// 	Debug(msg string)
}

type Application interface { // TODO
}

func NewServer(logger Logger, app Application) *Server {
	port := "8080"
	return &Server{
		httpServer: &http.Server{
			Addr: ":" + port,
			// Handler:        handler,
			MaxHeaderBytes: 1 << 20, // 1 MB
			ReadTimeout:    10 * time.Second,
			WriteTimeout:   10 * time.Second,
		},
	}
}

func (s *Server) Start(ctx context.Context) error {
	// TODO
	<-ctx.Done()
	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	// TODO
	return nil
}

// TODO
