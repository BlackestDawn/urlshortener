package main

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newSlogTestRouter(t *testing.T) (*gin.Engine, *bytes.Buffer) {
	t.Helper()
	gin.SetMode(gin.TestMode)

	var buf bytes.Buffer
	logger := slog.New(slog.NewJSONHandler(&buf, nil))

	router := gin.New()
	router.Use(requestid.New())
	router.Use(SlogLogger(logger))

	return router, &buf
}

func decodeLogLine(t *testing.T, buf *bytes.Buffer) map[string]any {
	t.Helper()
	var entry map[string]any
	require.NoError(t, json.Unmarshal(buf.Bytes(), &entry))
	return entry
}

func TestSlogLogger_LogsMethodPathAndStatus(t *testing.T) {
	router, buf := newSlogTestRouter(t)
	router.GET("/healthz", func(c *gin.Context) {
		c.String(http.StatusOK, "OK")
	})

	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	entry := decodeLogLine(t, buf)
	assert.Equal(t, "GET", entry["method"])
	assert.Equal(t, "/healthz", entry["path"])
	assert.Equal(t, float64(http.StatusOK), entry["status"])
}

func TestSlogLogger_LogsActualHandlerStatus_NotFound(t *testing.T) {
	router, buf := newSlogTestRouter(t)
	router.GET("/known", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})
	// no handler registered for /missing, so gin returns 404 itself

	req := httptest.NewRequest(http.MethodGet, "/missing", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	entry := decodeLogLine(t, buf)
	assert.Equal(t, float64(http.StatusNotFound), entry["status"])
}

func TestSlogLogger_IncludesQueryString(t *testing.T) {
	router, buf := newSlogTestRouter(t)
	router.GET("/search", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/search?q=abc", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	entry := decodeLogLine(t, buf)
	assert.Equal(t, "/search?q=abc", entry["path"])
}

func TestSlogLogger_RequestIDMatchesResponseHeader(t *testing.T) {
	router, buf := newSlogTestRouter(t)
	router.GET("/healthz", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	entry := decodeLogLine(t, buf)
	loggedID, ok := entry["request_id"].(string)
	require.True(t, ok, "request_id should be a string field")
	require.NotEmpty(t, loggedID)

	assert.Equal(t, w.Header().Get("X-Request-ID"), loggedID)
}

func TestSlogLogger_ClientProvidedRequestID_IsPreserved(t *testing.T) {
	router, buf := newSlogTestRouter(t)
	router.GET("/healthz", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	req.Header.Set("X-Request-ID", "client-supplied-id")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	entry := decodeLogLine(t, buf)
	assert.Equal(t, "client-supplied-id", entry["request_id"])
}
