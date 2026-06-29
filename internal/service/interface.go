package service

import "github.com/BlackestDawn/urlshortener/internal/domain"

type IShorten interface {
	Shorten(url string) (string, error)
	Resolve(code string) (string, error)
	GetStats(code string) (*domain.ShortUrl, error)
	Delete(code string) error
}
