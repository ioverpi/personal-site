package middleware

import (
	"log/slog"
	"time"

	"github.com/gin-gonic/gin"
)

// Logger middleware logs requests using slog
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		c.Next()

		latency := time.Since(start)
		status := c.Writer.Status()

		attrs := []any{
			"request_id", GetRequestID(c),
			"method", c.Request.Method,
			"path", path,
			"status", status,
			"latency_ms", float64(latency.Nanoseconds()) / 1e6,
			"ip", c.ClientIP(),
		}

		if query != "" {
			attrs = append(attrs, "query", query)
		}

		if len(c.Errors) > 0 {
			attrs = append(attrs, "errors", c.Errors.String())
		}

		switch {
		case status >= 500:
			slog.Error("request completed", attrs...)
		case status >= 400:
			slog.Warn("request completed", attrs...)
		default:
			slog.Info("request completed", attrs...)
		}
	}
}

// Log returns a logger with the request ID included
func Log(c *gin.Context) *slog.Logger {
	return slog.Default().With("request_id", GetRequestID(c))
}
