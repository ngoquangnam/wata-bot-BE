package middleware

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

type RequestLogMiddleware struct{}

func NewRequestLogMiddleware() *RequestLogMiddleware {
	return &RequestLogMiddleware{}
}

func (m *RequestLogMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Log request to terminal (using fmt.Printf to ensure it shows in console)
		fmt.Printf("\n━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
		fmt.Printf("📥 REQUEST [%s] %s", r.Method, r.URL.Path)
		if r.URL.RawQuery != "" {
			fmt.Printf("?%s", r.URL.RawQuery)
		}
		fmt.Printf("\n")
		fmt.Printf("   Origin: %s\n", r.Header.Get("Origin"))
		fmt.Printf("   User-Agent: %s\n", r.UserAgent())
		fmt.Printf("   Remote: %s\n", r.RemoteAddr)

		// Also log to logx for file logging
		logx.Infof("REQUEST [%s] %s %s - Origin: %s", r.Method, r.URL.Path, r.URL.RawQuery, r.Header.Get("Origin"))

		// Read and log request body if exists
		if r.Body != nil {
			requestBody, err := io.ReadAll(r.Body)
			if err == nil && len(requestBody) > 0 {
				// Restore body for handler to read
				r.Body = io.NopCloser(bytes.NewBuffer(requestBody))
				// Limit body size for logging (first 500 chars)
				bodyStr := string(requestBody)
				if len(bodyStr) > 500 {
					bodyStr = bodyStr[:500] + "... (truncated)"
				}
				fmt.Printf("   Body: %s\n", bodyStr)
				logx.Infof("Request Body: %s", bodyStr)
			} else if err == nil {
				// Body is empty, restore it anyway
				r.Body = io.NopCloser(bytes.NewBuffer(requestBody))
			}
		}

		// Create response writer wrapper to capture response
		rw := &responseLogWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
			body:           &bytes.Buffer{},
		}

		// Execute next handler
		next(rw, r)

		// Calculate duration
		duration := time.Since(start)

		// Log response to terminal
		fmt.Printf("📤 RESPONSE [%d] %s %s - %v\n", rw.statusCode, r.Method, r.URL.Path, duration)

		// Log response body if exists and not too large
		if rw.body.Len() > 0 {
			bodyStr := rw.body.String()
			if len(bodyStr) > 500 {
				bodyStr = bodyStr[:500] + "... (truncated)"
			}
			fmt.Printf("   Body: %s\n", bodyStr)
			logx.Infof("Response Body: %s", bodyStr)
		}

		fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n\n")

		// Also log to logx for file logging
		logx.Infof("RESPONSE [%d] %s %s - %v", rw.statusCode, r.Method, r.URL.Path, duration)
	}
}

type responseLogWriter struct {
	http.ResponseWriter
	statusCode int
	body       *bytes.Buffer
}

func (rw *responseLogWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func (rw *responseLogWriter) Write(b []byte) (int, error) {
	// Write to both the actual response and our buffer for logging
	rw.body.Write(b)
	return rw.ResponseWriter.Write(b)
}
