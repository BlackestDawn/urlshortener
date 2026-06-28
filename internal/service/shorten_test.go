package service

import (
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
