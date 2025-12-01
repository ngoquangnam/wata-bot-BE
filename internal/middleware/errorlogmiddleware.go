package middleware

import (
	"fmt"
	"net/http"
	"time"

	"wata-bot-BE/internal/utils"

	"github.com/zeromicro/go-zero/core/logx"
)

type ErrorLogMiddleware struct{}

func NewErrorLogMiddleware() *ErrorLogMiddleware {
	return &ErrorLogMiddleware{}
}

func (m *ErrorLogMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Create a response writer wrapper to capture status code
		rw := &responseWriter{
			ResponseWriter: w,
			statusCode:      http.StatusOK,
		}

		// Execute next handler
		next(rw, r)

		// Log errors (status code >= 400)
		if rw.statusCode >= 400 {
			duration := time.Since(start)

			context := map[string]interface{}{
				"method":      r.Method,
				"path":        r.URL.Path,
				"status":      rw.statusCode,
				"duration":    duration.String(),
				"remote_addr": r.RemoteAddr,
				"user_agent":  r.UserAgent(),
			}

			errorMsg := fmt.Sprintf("HTTP Error: %s %s", r.Method, r.URL.Path)
			utils.WriteErrorLogWithContext(errorMsg, nil, context)

			// Also log to go-zero logger
			logx.Errorf("HTTP Error [%d] %s %s - %v", rw.statusCode, r.Method, r.URL.Path, duration)
		}
	}
}

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}
