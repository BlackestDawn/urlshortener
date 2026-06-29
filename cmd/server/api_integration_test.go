package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/BlackestDawn/urlshortener/api"
	"github.com/BlackestDawn/urlshortener/internal/domain"
	"github.com/BlackestDawn/urlshortener/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const testDomain = "short.test"

func newTestRouter(t *testing.T) (*gin.Engine, *service.MockIShorten) {
	t.Helper()
	gin.SetMode(gin.TestMode)

	srv := service.NewMockIShorten(t)
	api := NewApiController(srv, testDomain)

	router := gin.New()
	router.Use(ErrorHandler())

	router.GET("/healthz", api.GetHealth)
	router.GET("/:code", api.Redirect)

	apiRoute := router.Group("/api")
	{
		v1Route := apiRoute.Group("/v1")
		{
			linksRoute := v1Route.Group("/links")
			{
				linksRoute.GET("/:code/stats", api.GetStats)
				linksRoute.GET("/:code", api.GetSingle)
				linksRoute.POST("", api.Create)
				linksRoute.DELETE("/:code", api.Remove)
			}
		}
	}

	return router, srv
}

func doRequest(t *testing.T, router *gin.Engine, method, target string, body any) *httptest.ResponseRecorder {
	t.Helper()

	var reqBody *bytes.Buffer
	if body != nil {
		raw, err := json.Marshal(body)
		require.NoError(t, err)
		reqBody = bytes.NewBuffer(raw)
	} else {
		reqBody = bytes.NewBuffer(nil)
	}

	req := httptest.NewRequest(method, target, reqBody)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	return rec
}

func TestHealthz_ReturnsOK(t *testing.T) {
	router, _ := newTestRouter(t)

	rec := doRequest(t, router, http.MethodGet, "/healthz", nil)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), "OK")
}

func TestCreate_ValidUrl_ReturnsCreatedWithShortenedUrl(t *testing.T) {
	router, srv := newTestRouter(t)
	srv.EXPECT().
		Shorten("https://example.com").
		Return("abc123", nil)

	rec := doRequest(t, router, http.MethodPost, "/api/v1/links", api.UrlDto{Url: "https://example.com"})

	require.Equal(t, http.StatusCreated, rec.Code)

	var resp api.ShortenedUrlDto
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
	assert.Equal(t, "https://short.test/abc123", resp.ShortenedUrl)
}

func TestCreate_InvalidUrl_ReturnsBadRequest(t *testing.T) {
	router, srv := newTestRouter(t)
	srv.EXPECT().
		Shorten("not-a-url").
		Return("", domain.ErrInvalidUrl)

	rec := doRequest(t, router, http.MethodPost, "/api/v1/links", api.UrlDto{Url: "not-a-url"})

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestCreate_MalformedJson_ReturnsUnprocessableEntity(t *testing.T) {
	router, _ := newTestRouter(t)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/links", bytes.NewBufferString("{not-json"))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)
}

func TestGetSingle_KnownCode_ReturnsUrl(t *testing.T) {
	router, srv := newTestRouter(t)
	srv.EXPECT().
		Resolve("abc123").
		Return("https://example.com", nil)

	rec := doRequest(t, router, http.MethodGet, "/api/v1/links/abc123", nil)

	require.Equal(t, http.StatusOK, rec.Code)

	var resp api.UrlDto
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
	assert.Equal(t, "https://example.com", resp.Url)
}

func TestGetSingle_UnknownCode_ReturnsNotFound(t *testing.T) {
	router, srv := newTestRouter(t)
	srv.EXPECT().
		Resolve("missing").
		Return("", domain.ErrNotFound)

	rec := doRequest(t, router, http.MethodGet, "/api/v1/links/missing", nil)

	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestGetStats_KnownCode_ReturnsStats(t *testing.T) {
	router, srv := newTestRouter(t)
	created := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	srv.EXPECT().
		GetStats("abc123").
		Return(&domain.ShortUrl{
			Code:        "abc123",
			OriginalUrl: "https://example.com",
			Clicks:      5,
			CreatedAt:   created,
		}, nil)

	rec := doRequest(t, router, http.MethodGet, "/api/v1/links/abc123/stats", nil)

	require.Equal(t, http.StatusOK, rec.Code)

	var resp api.UrlStatsDto
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
	assert.Equal(t, "abc123", resp.Code)
	assert.Equal(t, "https://example.com", resp.Url)
	assert.Equal(t, 5, resp.Hits)
	assert.True(t, created.Equal(resp.CreatedAt))
}

func TestGetStats_UnknownCode_ReturnsNotFound(t *testing.T) {
	router, srv := newTestRouter(t)
	srv.EXPECT().
		GetStats("missing").
		Return(nil, domain.ErrNotFound)

	rec := doRequest(t, router, http.MethodGet, "/api/v1/links/missing/stats", nil)

	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestRemove_KnownCode_ReturnsNoContent(t *testing.T) {
	router, srv := newTestRouter(t)
	srv.EXPECT().
		Delete("abc123").
		Return(nil)

	rec := doRequest(t, router, http.MethodDelete, "/api/v1/links/abc123", nil)

	assert.Equal(t, http.StatusNoContent, rec.Code)
}

func TestRemove_UnknownCode_ReturnsNotFound(t *testing.T) {
	router, srv := newTestRouter(t)
	srv.EXPECT().
		Delete("missing").
		Return(domain.ErrNotFound)

	rec := doRequest(t, router, http.MethodDelete, "/api/v1/links/missing", nil)

	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestRedirect_KnownCode_ReturnsPermanentRedirect(t *testing.T) {
	router, srv := newTestRouter(t)
	srv.EXPECT().
		Resolve("abc123").
		Return("https://example.com", nil)

	rec := doRequest(t, router, http.MethodGet, "/abc123", nil)

	assert.Equal(t, http.StatusPermanentRedirect, rec.Code)
	assert.Equal(t, "https://example.com", rec.Header().Get("Location"))
}

func TestRedirect_UnknownCode_ReturnsNotFound(t *testing.T) {
	router, srv := newTestRouter(t)
	srv.EXPECT().
		Resolve("missing").
		Return("", domain.ErrNotFound)

	rec := doRequest(t, router, http.MethodGet, "/missing", nil)

	assert.Equal(t, http.StatusNotFound, rec.Code)
}
