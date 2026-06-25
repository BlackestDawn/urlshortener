//go:build integration

package repository

import (
	"testing"

	"github.com/BlackestDawn/urlshortener/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRepository_Create_PersistsRecord(t *testing.T) {
	result, err := testRepo.Create("https://example.com")
	require.NoError(t, err)
	assert.NotEmpty(t, result.Code)
	assert.Equal(t, "https://example.com", result.OriginalUrl)
}

func TestRepository_FindByCode_ReturnsRecord(t *testing.T) {
	created, err := testRepo.Create("https://example.com/find")
	require.NoError(t, err)

	found, err := testRepo.FindByCode(created.Code)
	require.NoError(t, err)
	assert.Equal(t, created.Code, found.Code)
}

func TestRepository_FindByCode_ReturnsErrNotFound(t *testing.T) {
	_, err := testRepo.FindByCode("1234567890abcdef")
	require.ErrorIs(t, err, domain.ErrNotFound)
}
