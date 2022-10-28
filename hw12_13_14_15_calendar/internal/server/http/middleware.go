package http

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

type LoggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func newLoggingResponseWriter(w http.ResponseWriter) *LoggingResponseWriter {
	return &LoggingResponseWriter{w, http.StatusOK}
}

func LoggingMiddleware(logg Logger) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			startTime := time.Now()
			writter := newLoggingResponseWriter(w)
			next.ServeHTTP(writter, r)
			latency := time.Since(startTime)
			logg.Info(fmt.Sprintf("%s %s %s %s %s %d %d ms %s",
				r.RemoteAddr,
				fmt.Sprint("[", startTime.Format("02/Jan/2006 15:04:05 -0700"), "]"),
				r.Method,
				r.RequestURI,
				r.Proto,
				writter.statusCode,
				latency.Microseconds(),
				r.UserAgent()))
		})
	}
}
