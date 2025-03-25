package middleware

import (
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

type ResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *ResponseWriter) WriteHeader(statusCode int) {
	rw.statusCode = statusCode
	rw.ResponseWriter.WriteHeader(statusCode)
}

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		responseWriter := &ResponseWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}

		logger := log.With().
			Str("request_id", uuid.New().String()).
			Str("method", r.Method).
			Str("path", r.URL.Path).
			Str("remote_addr", r.RemoteAddr).
			Str("user_agent", r.UserAgent()).
			Logger()

		ctx := logger.WithContext(r.Context())
		r = r.WithContext(ctx)

		next.ServeHTTP(responseWriter, r)

		responseTime := time.Since(start)

		logEvent := logger.Info()
		if responseWriter.statusCode >= 400 && responseWriter.statusCode < 500 {
			logEvent = logger.Warn()
		} else if responseWriter.statusCode >= 500 {
			logEvent = logger.Error()
		}

		logEvent.
			Int("status", responseWriter.statusCode).
			Dur("duration_ms", responseTime).
			Msg("request completed")
	})
}
