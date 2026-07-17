package main

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewTimeoutHandler_HandlerFinishesInTime_ReturnsNormally(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/fast", func(c *gin.Context) {
		c.String(http.StatusOK, "done")
	})

	handler := newTimeoutHandler(router, 200*time.Millisecond)

	req := httptest.NewRequest(http.MethodGet, "/fast", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "done", w.Body.String())
}

func TestNewTimeoutHandler_HandlerTooSlow_Returns503(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/slow", func(c *gin.Context) {
		time.Sleep(100 * time.Millisecond)
		c.String(http.StatusOK, "done")
	})

	handler := newTimeoutHandler(router, 10*time.Millisecond)

	req := httptest.NewRequest(http.MethodGet, "/slow", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	require.Equal(t, http.StatusServiceUnavailable, w.Code)
	assert.Equal(t, "request timed out", w.Body.String())
}

func TestNewTimeoutHandler_CancelsHandlerContext(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	ctxErr := make(chan error, 1)
	router.GET("/slow", func(c *gin.Context) {
		<-c.Request.Context().Done()
		ctxErr <- c.Request.Context().Err()
	})

	handler := newTimeoutHandler(router, 10*time.Millisecond)

	req := httptest.NewRequest(http.MethodGet, "/slow", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	require.Equal(t, http.StatusServiceUnavailable, w.Code)

	select {
	case err := <-ctxErr:
		assert.ErrorIs(t, err, context.DeadlineExceeded)
	case <-time.After(time.Second):
		t.Fatal("handler's context was never cancelled")
	}
}
