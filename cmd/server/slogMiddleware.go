package main

import (
	"log/slog"
	"time"

	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
)

// SlogLogger logs each request as a single structured record, tagged with
// the request ID assigned by requestid.New() so log lines can be correlated
// with that request end to end. Must run after requestid.New() in the chain.
func SlogLogger(logger *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		if raw := c.Request.URL.RawQuery; raw != "" {
			path += "?" + raw
		}

		c.Next()

		logger.Info("request",
			slog.String("request_id", requestid.Get(c)),
			slog.String("method", c.Request.Method),
			slog.String("path", path),
			slog.Int("status", c.Writer.Status()),
			slog.Duration("latency", time.Since(start)),
			slog.String("client_ip", c.ClientIP()),
		)
	}
}
