package repository

import "github.com/BlackestDawn/urlshortener/internal/domain"

func entryToDomain(entry ShortUrl) *domain.ShortUrl {
	return &domain.ShortUrl{
		ID:          entry.ID,
		CreatedAt:   entry.CreatedAt,
		Code:        entry.Code,
		OriginalUrl: entry.OriginalUrl,
		Clicks:      int(entry.Clicks),
	}
}
