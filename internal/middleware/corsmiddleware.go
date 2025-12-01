package middleware

import (
	"net/http"
	"strings"

	"github.com/zeromicro/go-zero/core/logx"
)

type CorsMiddleware struct {
	allowedOrigins []string
}

func NewCorsMiddleware() *CorsMiddleware {
	return &CorsMiddleware{
		allowedOrigins: []string{}, // Empty means allow all origins
	}
}

// NewCorsMiddlewareWithOrigins creates CORS middleware with specific allowed origins
func NewCorsMiddlewareWithOrigins(origins []string) *CorsMiddleware {
	return &CorsMiddleware{
		allowedOrigins: origins,
	}
}

func (m *CorsMiddleware) isOriginAllowed(origin string) bool {
	// If no origins specified, allow all
	if len(m.allowedOrigins) == 0 {
		return true
	}

	// Normalize origin (remove trailing slash for comparison)
	normalizedOrigin := strings.TrimSuffix(origin, "/")

	// Check if origin is in allowed list
	for _, allowed := range m.allowedOrigins {
		if allowed == "*" {
			return true
		}
		// Normalize allowed origin too
		normalizedAllowed := strings.TrimSuffix(allowed, "/")
		if normalizedAllowed == normalizedOrigin {
			return true
		}
	}
	return false
}

func (m *CorsMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")

		// Debug logging
		logx.Infof("CORS Request - Origin: %s, Method: %s, Path: %s", origin, r.Method, r.URL.Path)

		// Create a response writer wrapper to ensure headers are set
		corsWriter := &corsResponseWriter{
			ResponseWriter: w,
			allowedOrigins: m.allowedOrigins,
			origin:         origin,
		}

		// Handle preflight OPTIONS request
		if r.Method == http.MethodOptions {
			logx.Infof("CORS: Handling preflight OPTIONS request")
			corsWriter.setCorsHeaders()
			w.WriteHeader(http.StatusNoContent)
			return
		}

		// Continue with the next handler
		next(corsWriter, r)
	}
}

// corsResponseWriter wraps http.ResponseWriter to ensure CORS headers are always set
type corsResponseWriter struct {
	http.ResponseWriter
	allowedOrigins []string
	origin         string
	headersSet     bool
}

func (rw *corsResponseWriter) setCorsHeaders() {
	if rw.headersSet {
		return
	}
	rw.headersSet = true

	// Always set CORS headers for all requests
	rw.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS, PATCH")
	rw.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With, Accept, Origin")
	rw.Header().Set("Access-Control-Max-Age", "3600")
	rw.Header().Set("Access-Control-Expose-Headers", "Content-Length, Content-Type")

	// Set Access-Control-Allow-Origin based on origin
	if rw.origin != "" {
		normalizedOrigin := strings.TrimSuffix(rw.origin, "/")
		allowed := false

		if len(rw.allowedOrigins) == 0 {
			allowed = true
		} else {
			for _, allowedOrigin := range rw.allowedOrigins {
				normalizedAllowed := strings.TrimSuffix(allowedOrigin, "/")
				if normalizedAllowed == normalizedOrigin || allowedOrigin == "*" {
					allowed = true
					break
				}
			}
		}

		if allowed {
			rw.Header().Set("Access-Control-Allow-Origin", rw.origin)
			rw.Header().Set("Access-Control-Allow-Credentials", "true")
			logx.Infof("CORS: Allowing origin %s", rw.origin)
		} else {
			logx.Errorf("CORS: Origin %s not allowed. Allowed origins: %v", rw.origin, rw.allowedOrigins)
		}
	} else {
		// No origin header
		if len(rw.allowedOrigins) == 0 {
			rw.Header().Set("Access-Control-Allow-Origin", "*")
		}
	}
}

func (rw *corsResponseWriter) WriteHeader(code int) {
	rw.setCorsHeaders()
	rw.ResponseWriter.WriteHeader(code)
}

func (rw *corsResponseWriter) Write(b []byte) (int, error) {
	rw.setCorsHeaders()
	return rw.ResponseWriter.Write(b)
}
