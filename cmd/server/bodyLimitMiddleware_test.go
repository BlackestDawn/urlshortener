package main

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/BlackestDawn/urlshortener/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func newBodyLimitRouter(limit int64) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(MaxBodySize(limit))
	router.POST("/echo", func(c *gin.Context) {
		body, err := io.ReadAll(c.Request.Body)
		if err != nil {
			c.Error(err)
			return
		}
		c.String(http.StatusOK, "%d", len(body))
	})
	return router
}

func TestMaxBodySize_UnderLimit_PassesThrough(t *testing.T) {
	router := newBodyLimitRouter(16)

	req := httptest.NewRequest(http.MethodPost, "/echo", bytes.NewBufferString("short"))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "5", w.Body.String())
}

func TestMaxBodySize_ExactlyAtLimit_PassesThrough(t *testing.T) {
	router := newBodyLimitRouter(5)

	req := httptest.NewRequest(http.MethodPost, "/echo", bytes.NewBufferString("12345"))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "5", w.Body.String())
}

func TestMaxBodySize_OverLimit_ReadFailsWithMaxBytesError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(MaxBodySize(5))
	router.Use(ErrorHandler())
	router.POST("/echo", func(c *gin.Context) {
		body, err := io.ReadAll(c.Request.Body)
		if err != nil {
			c.Error(err)
			return
		}
		c.String(http.StatusOK, "%d", len(body))
	})

	req := httptest.NewRequest(http.MethodPost, "/echo", bytes.NewBufferString("123456"))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestMaxBodySize_OverLimit_ErrorTypeIsMaxBytesError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(MaxBodySize(5))

	var readErr error
	router.POST("/echo", func(c *gin.Context) {
		_, readErr = io.ReadAll(c.Request.Body)
		c.Status(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodPost, "/echo", bytes.NewBufferString("123456"))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	var maxBytesErr *http.MaxBytesError
	require.True(t, errors.As(readErr, &maxBytesErr), "expected *http.MaxBytesError, got %T: %v", readErr, readErr)
}

func TestMaxBodySize_OversizedCreateRequest_Returns413(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(MaxBodySize(16))
	router.Use(ErrorHandler())
	api := NewApiController(nil, "example.com")
	router.POST("/links", api.Create)

	oversized := `{"url":"` + strings.Repeat("a", 100) + `"}`
	req := httptest.NewRequest(http.MethodPost, "/links", bytes.NewBufferString(oversized))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusRequestEntityTooLarge, w.Code)
	assert.Contains(t, w.Body.String(), "too large")
}

func TestMaxBodySize_NormalSizedCreateRequest_ReachesHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(MaxBodySize(maxRequestBodyBytes))
	router.Use(ErrorHandler())
	srv := service.NewMockIShorten(t)
	srv.EXPECT().Shorten(mock.Anything, "https://example.com").Return("abc123", nil)
	api := NewApiController(srv, "example.com")
	router.POST("/links", api.Create)

	req := httptest.NewRequest(http.MethodPost, "/links", bytes.NewBufferString(`{"url":"https://example.com"}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusCreated, w.Code)
}
