package internalhttp

import (
	"fmt"
	"net/http"
	"time"

	appLogger "github.com/rodiond26/go-2022-05/hw12_13_14_15_calendar/internal/logger"
)

func loggingMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		h.ServeHTTP(w, r)
		latency := time.Since(startTime)

		l, _ := appLogger.NewLogger("dev", "info")
		l.Info(fmt.Sprintf("%s [%s] %s %s %s %d %s \"%s\"",
			r.RemoteAddr,
			time.Now().Format("2006-01-02 15:04:05 -0700"),
			r.Method,
			r.URL.Path,
			r.Proto,
			http.StatusOK,
			latency,
			r.UserAgent(),
		))
	})
}
