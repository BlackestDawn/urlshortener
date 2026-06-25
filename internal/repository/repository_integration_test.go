//go:build integration

package repository

import (
	"testing"

	"github.com/BlackestDawn/urlshortener/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func resetDB(t *testing.T) {
	t.Helper()
	_, err := testDB.Exec("TRUNCATE TABLE short_urls")
	require.NoError(t, err)
}

func createTestRecord(t *testing.T, url string) *domain.ShortUrl {
	t.Helper()
	created, err := testRepo.Create(url)
	require.NoError(t, err)
	return created
}

func TestRepository_Create_PersistsRecord(t *testing.T) {
	resetDB(t)

	result, err := testRepo.Create("https://example.com")
	require.NoError(t, err)
	assert.NotEmpty(t, result.Code)
	assert.Equal(t, "https://example.com", result.OriginalUrl)
}

func TestRepository_FindByCode_ReturnsRecord(t *testing.T) {
	resetDB(t)

	created := createTestRecord(t, "https://example.com/find")

	found, err := testRepo.FindByCode(created.Code)
	require.NoError(t, err)
	assert.Equal(t, created.Code, found.Code)
}

func TestRepository_FindByCode_ReturnsErrNotFound(t *testing.T) {
	resetDB(t)

	_, err := testRepo.FindByCode("1234567890abcdef")
	require.ErrorIs(t, err, domain.ErrNotFound)
}

func TestRepository_IncrementClicks_CounterIncreases(t *testing.T) {
	resetDB(t)

	created := createTestRecord(t, "https://example.com/clicks")

	err := testRepo.IncrementClicks(created.Code)
	require.NoError(t, err)

	found, err := testRepo.FindByCode(created.Code)
	require.NoError(t, err)
	assert.Equal(t, created.Clicks+1, found.Clicks)
}

func TestRepository_IncrementClicks_NonExistentCode(t *testing.T) {
	resetDB(t)

	err := testRepo.IncrementClicks("1234567890abcdef")
	require.ErrorIs(t, err, domain.ErrNotFound)
}

func TestRepository_List_ReturnsAllRecords(t *testing.T) {
	resetDB(t)

	urls := []string{
		"https://example.com/1",
		"https://example.com/2",
		"https://example.com/3",
	}
	for _, u := range urls {
		createTestRecord(t, u)
	}

	results, total, err := testRepo.List(1, 10, "")
	require.NoError(t, err)
	assert.Equal(t, len(urls), total)
	assert.Len(t, results, len(urls))
}

func TestRepository_List_EmptyDB(t *testing.T) {
	resetDB(t)

	results, total, err := testRepo.List(1, 10, "")
	require.NoError(t, err)
	assert.Equal(t, 0, total)
	assert.Empty(t, results)
}

func TestRepository_Delete_RemovesRecord(t *testing.T) {
	resetDB(t)

	created := createTestRecord(t, "https://example.com/")

	err := testRepo.Delete(created.Code)
	require.NoError(t, err)

	_, err = testRepo.FindByCode(created.Code)
	assert.ErrorIs(t, err, domain.ErrNotFound)
}

func TestRepository_Delete_NonExistentCode(t *testing.T) {
	resetDB(t)

	err := testRepo.Delete("1234567890abcdef")
	assert.NoError(t, err)
}

func TestRepository_Create_DuplicateCode_ReturnsError(t *testing.T) {
	resetDB(t)

	url := "https://example.com"

	createTestRecord(t, url)

	_, err := testRepo.Create(url)
	assert.Error(t, err)
}
