package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// MaxBodySize rejects request bodies larger than limit bytes.
// Reads past the limit fail with a *http.MaxBytesError, which ErrorHandler
// translates into a 413 response.
func MaxBodySize(limit int64) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, limit)
		c.Next()
	}
}
