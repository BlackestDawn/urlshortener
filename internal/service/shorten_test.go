package service

import (
	"database/sql"
	"testing"

	"github.com/BlackestDawn/urlshortener/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const validUrl = "https://example.com"
const validCode = "abc123"

func TestShorten_ValidUrl_ReturnsCode(t *testing.T) {
	repo := domain.NewMockIRepository(t)
	repo.EXPECT().
		Create(validUrl).
		Return(&domain.ShortUrl{Code: validCode}, nil)

	svc := NewShortenService(repo)

	code, err := svc.Shorten(validUrl)
	require.NoError(t, err)
	assert.Equal(t, validCode, code)
}

func TestShorten_CallsRepoCreate_Once(t *testing.T) {
	repo := domain.NewMockIRepository(t)
	repo.EXPECT().
		Create(validUrl).
		Return(&domain.ShortUrl{Code: validCode}, nil).
		Once()

	svc := NewShortenService(repo)
	_, err := svc.Shorten(validUrl)
	assert.NoError(t, err)
}

func TestShorten_InvalidURL_ReturnsError(t *testing.T) {
	repo := domain.NewMockIRepository(t)

	svc := NewShortenService(repo)

	_, err := svc.Shorten("invalid_url")
	assert.ErrorIs(t, domain.ErrInvalidUrl, err)
}

func TestShorten_RepoError_PropagatesError(t *testing.T) {
	repo := domain.NewMockIRepository(t)
	repo.EXPECT().
		Create(validUrl).
		Return(nil, domain.ErrNotFound)

	svc := NewShortenService(repo)

	_, err := svc.Shorten(validUrl)
	assert.ErrorIs(t, domain.ErrNotFound, err)
}

func TestResolve_KnownCode_ReturnsURL(t *testing.T) {
	repo := domain.NewMockIRepository(t)
	repo.EXPECT().
		FindByCode(validCode).
		Return(&domain.ShortUrl{OriginalUrl: validUrl}, nil)
	repo.EXPECT().
		IncrementClicks(validCode).
		Return(nil)

	svc := NewShortenService(repo)

	url, err := svc.Resolve(validCode)
	assert.NoError(t, err)
	assert.Equal(t, validUrl, url)
}

func TestResolve_UnknownCode_ReturnsErrNotFound(t *testing.T) {
	repo := domain.NewMockIRepository(t)
	repo.EXPECT().
		FindByCode("invalid_code").
		Return(nil, sql.ErrNoRows)

	svc := NewShortenService(repo)

	_, err := svc.Resolve("invalid_code")
	assert.ErrorIs(t, domain.ErrNotFound, err)
}
func TestResolve_IncrementsClicks(t *testing.T) {
	repo := domain.NewMockIRepository(t)
	repo.EXPECT().
		FindByCode(validCode).
		Return(&domain.ShortUrl{OriginalUrl: validUrl}, nil)
	repo.EXPECT().
		IncrementClicks(validCode).
		Return(nil).
		Once()

	svc := NewShortenService(repo)

	_, err := svc.Resolve(validCode)
	assert.NoError(t, err)
}
