package service

import (
	"github.com/BlackestDawn/urlshortener/internal/domain"
)

type ShortenService struct {
	repo domain.IRepository
}

func NewShortenService(repo domain.IRepository) *ShortenService {
	return &ShortenService{repo: repo}
}

func (s *ShortenService) Shorten(url string) (string, error) {
	if ret, _ := domain.ValidateURL(url); !ret {
		return "", domain.ErrInvalidUrl
	}

	entry, err := s.repo.Create(url)
	if err != nil {
		return "", err
	}

	return entry.Code, nil
}

func (s *ShortenService) Resolve(code string) (string, error) {
	entry, err := s.repo.FindByCode(code)
	if err != nil {
		return "", err
	}

	err = s.repo.IncrementClicks(code)
	if err != nil {
		return "", err
	}

	return entry.OriginalUrl, nil
}
