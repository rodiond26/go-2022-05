package internalhttp

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
)

type Server struct {
	httpServer *http.Server
	logger     Logger
	router     *httprouter.Router
}

type Logger interface {
	Info(msg string)
	Error(msg string)
	Warn(msg string)
	Debug(msg string)
}

type Application interface { // TODO
}

// func NewServer(logger Logger, app Application) *Server {
// 	port := "8080"
// 	return &Server{
// 		httpServer: &http.Server{
// 			Addr: ":" + port,
// 			// Handler:        handler,
// 			MaxHeaderBytes: 1 << 20, // 1 MB
// 			ReadTimeout:    10 * time.Second,
// 			WriteTimeout:   10 * time.Second,
// 		},
// 	}
// }

func NewServer(logger Logger) *Server { // app Application
	serv := &Server{logger: logger}
	serv.router = httprouter.New()

	serv.router.GET("/", func(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
		text := "Hello world!"
		fmt.Fprint(writer, text)
	})

	return serv
}

func (s *Server) Start(ctx context.Context, addr string) error {
	s.logger.Info("HTTP server [" + addr + "] starting...")
	s.httpServer = &http.Server{
		Addr:         addr,
		Handler:      s.router,
		WriteTimeout: 5 * time.Second,
		ReadTimeout:  5 * time.Second,
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
	s.logger.Info("HTTP server was stopped...")
	if err := s.httpServer.Shutdown(ctx); err != nil {
		return err
	}
	return nil
}
