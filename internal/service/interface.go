package service

import (
	"context"

	"github.com/BlackestDawn/urlshortener/internal/domain"
)

type IShorten interface {
	Shorten(ctx context.Context, url string) (string, error)
	Resolve(ctx context.Context, code string) (string, error)
	GetStats(ctx context.Context, code string) (*domain.ShortUrl, error)
	Delete(ctx context.Context, code string) error
}
