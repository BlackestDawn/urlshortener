package api

import "github.com/BlackestDawn/urlshortener/internal/domain"

func EntityToStatDto(entity *domain.ShortUrl) *UrlStatsDto {
	return &UrlStatsDto{
		Url:       entity.OriginalUrl,
		Code:      entity.Code,
		Hits:      entity.Clicks,
		CreatedAt: entity.CreatedAt,
	}
}
