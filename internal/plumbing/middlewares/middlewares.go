package middlewares

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type Middleware func(http.HandlerFunc) http.HandlerFunc

// Logging logs all requests with its path and the time it took to process
func Logging() Middleware {
	// Create a new Middleware
	return func(f http.HandlerFunc) http.HandlerFunc {
		// Define the http.HandlerFunc
		return func(w http.ResponseWriter, r *http.Request) {
			requestId := uuid.New()

			originalLogger := slog.Default()
			requestLogger := slog.With("requestId", requestId)
			slog.SetDefault(requestLogger)

			start := time.Now()
			slog.Info("Start request", "path", r.URL.Path, "start_time", start)

			defer func() {
				slog.Info("End request", "path", r.URL.Path, "duration", time.Since(start))
				// reinstate the logger
				slog.SetDefault(originalLogger)
			}()

			// Call the next middleware/handler in chain
			f(w, r)
		}
	}
}

// Chain applies middlewares to a http.HandlerFunc
func Chain(f http.HandlerFunc, middlewares ...Middleware) http.HandlerFunc {
	for _, m := range middlewares {
		f = m(f)
	}
	return f
}
