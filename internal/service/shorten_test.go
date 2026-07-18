package service

import (
	"testing"

	"github.com/BlackestDawn/urlshortener/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

const validUrl = "https://example.com"
const validCode = "abc123"

func TestShorten_ValidUrl_ReturnsCode(t *testing.T) {
	repo := domain.NewMockIRepository(t)
	repo.EXPECT().
		Create(mock.Anything, validUrl).
		Return(&domain.ShortUrl{Code: validCode}, nil)

	svc := NewShortenService(repo)

	code, err := svc.Shorten(t.Context(), validUrl)
	require.NoError(t, err)
	assert.Equal(t, validCode, code)
}

func TestShorten_CallsRepoCreate_Once(t *testing.T) {
	repo := domain.NewMockIRepository(t)
	repo.EXPECT().
		Create(mock.Anything, validUrl).
		Return(&domain.ShortUrl{Code: validCode}, nil).
		Once()

	svc := NewShortenService(repo)
	_, err := svc.Shorten(t.Context(), validUrl)
	assert.NoError(t, err)
}

func TestShorten_InvalidURL_ReturnsError(t *testing.T) {
	repo := domain.NewMockIRepository(t)

	svc := NewShortenService(repo)

	_, err := svc.Shorten(t.Context(), "invalid_url")
	assert.ErrorIs(t, err, domain.ErrInvalidUrl)
}

func TestShorten_RepoError_PropagatesError(t *testing.T) {
	repo := domain.NewMockIRepository(t)
	repo.EXPECT().
		Create(mock.Anything, validUrl).
		Return(nil, domain.ErrNotFound)

	svc := NewShortenService(repo)

	_, err := svc.Shorten(t.Context(), validUrl)
	assert.ErrorIs(t, err, domain.ErrNotFound)
}

func TestResolve_KnownCode_ReturnsURL(t *testing.T) {
	repo := domain.NewMockIRepository(t)
	repo.EXPECT().
		IncrementClicks(mock.Anything, validCode).
		Return(&domain.ShortUrl{OriginalUrl: validUrl}, nil)

	svc := NewShortenService(repo)

	url, err := svc.Resolve(t.Context(), validCode)
	assert.NoError(t, err)
	assert.Equal(t, validUrl, url)
}

func TestResolve_UnknownCode_ReturnsErrNotFound(t *testing.T) {
	repo := domain.NewMockIRepository(t)
	repo.EXPECT().
		IncrementClicks(mock.Anything, "invalid_code").
		Return(nil, domain.ErrNotFound)

	svc := NewShortenService(repo)

	_, err := svc.Resolve(t.Context(), "invalid_code")
	assert.ErrorIs(t, err, domain.ErrNotFound)
}
func TestResolve_IncrementsClicks(t *testing.T) {
	repo := domain.NewMockIRepository(t)
	repo.EXPECT().
		IncrementClicks(mock.Anything, validCode).
		Return(&domain.ShortUrl{OriginalUrl: validUrl}, nil).
		Once()

	svc := NewShortenService(repo)

	_, err := svc.Resolve(t.Context(), validCode)
	assert.NoError(t, err)
}

func TestGetStats_KnownCode_ReturnsStruct(t *testing.T) {
	repo := domain.NewMockIRepository(t)
	repo.EXPECT().
		FindByCode(mock.Anything, validCode).
		Return(&domain.ShortUrl{OriginalUrl: validUrl, Code: validCode, Clicks: 10}, nil)

	svc := NewShortenService(repo)

	entry, err := svc.GetStats(t.Context(), validCode)
	assert.NoError(t, err)
	assert.Equal(t, validCode, entry.Code)
	assert.Equal(t, validUrl, entry.OriginalUrl)
	assert.Equal(t, 10, entry.Clicks)
}
func TestGetStats_UnknownCode_ReturnsErrNotFound(t *testing.T) {
	repo := domain.NewMockIRepository(t)
	repo.EXPECT().
		FindByCode(mock.Anything, "invalid_code").
		Return(nil, domain.ErrNotFound)

	svc := NewShortenService(repo)

	_, err := svc.GetStats(t.Context(), "invalid_code")
	assert.ErrorIs(t, err, domain.ErrNotFound)
}
func TestDelete_KnownCode_CallsRepo(t *testing.T) {
	repo := domain.NewMockIRepository(t)
	repo.EXPECT().
		Delete(mock.Anything, validCode).
		Return(nil).
		Once()

	svc := NewShortenService(repo)

	err := svc.Delete(t.Context(), validCode)
	assert.NoError(t, err)
}
