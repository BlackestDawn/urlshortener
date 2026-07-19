package service

import (
	"context"

	"github.com/BlackestDawn/urlshortener/internal/domain"
)

type ShortenService struct {
	repo domain.IRepository
}

func NewShortenService(repo domain.IRepository) *ShortenService {
	return &ShortenService{repo: repo}
}

func (s *ShortenService) Shorten(ctx context.Context, url string) (string, error) {
	if ret, _ := domain.ValidateURL(url); !ret {
		return "", domain.ErrInvalidUrl
	}

	entry, err := s.repo.Create(ctx, url)
	if err != nil {
		return "", err
	}

	return entry.Code, nil
}

func (s *ShortenService) Resolve(ctx context.Context, code string) (string, error) {
	entry, err := s.repo.IncrementClicks(ctx, code)
	if err != nil {
		return "", err
	}

	return entry.OriginalUrl, nil
}

func (s *ShortenService) GetStats(ctx context.Context, code string) (*domain.ShortUrl, error) {
	entry, err := s.repo.FindByCode(ctx, code)
	if err != nil {
		return nil, err
	}

	return entry, nil
}

func (s *ShortenService) Delete(ctx context.Context, code string) error {
	return s.repo.Delete(ctx, code)
}
