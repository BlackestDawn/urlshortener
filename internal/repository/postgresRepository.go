package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/BlackestDawn/urlshortener/config"
	"github.com/BlackestDawn/urlshortener/internal/domain"
)

type PostgresRepository struct {
	QBQueries *Queries
}

func NewPGRepository(cfg config.Config) (*PostgresRepository, error) {
	repo := new(PostgresRepository)

	db, err := sql.Open("postgres", cfg.DBUrl)
	if err != nil {
		return nil, err
	}

	repo.QBQueries = New(db)

	return repo, nil
}

func (r *PostgresRepository) CreateShortUrl(url string) (*domain.ShortUrl, error) {
	if res, err := domain.ValidateURL(url); !res {
		return nil, fmt.Errorf("Invalid URL: %w", err)
	}

	code, err := domain.GenerateCode(url)
	if err != nil {
		return nil, fmt.Errorf("Error generating code: %w", err)
	}

	entry, err := r.QBQueries.CreateShortUrl(context.Background(), CreateShortUrlParams{
		Code:        code,
		OriginalUrl: url,
	})
	if err != nil {
		return nil, err
	}

	return entryToDomain(entry), nil
}

func (r *PostgresRepository) GetByCode(code string) (*domain.ShortUrl, error) {
	entry, err := r.QBQueries.GetByCode(context.Background(), code)
	if err != nil {
		return nil, err
	}

	return entryToDomain(entry), nil
}

func (r *PostgresRepository) Remove(code string) error {
	return r.QBQueries.DeleteByCode(context.Background(), code)
}
