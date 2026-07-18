package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/BlackestDawn/urlshortener/config"
	"github.com/BlackestDawn/urlshortener/internal/domain"
	_ "github.com/jackc/pgx/v5/stdlib"
)

const (
	maxOpenConns    = 25
	maxIdleConns    = 25
	connMaxLifetime = 5 * time.Minute
	connMaxIdleTime = 5 * time.Minute
)

type PostgresRepository struct {
	QBQueries *Queries
}

func NewPGRepository(cfg *config.Config) (*PostgresRepository, error) {
	repo := new(PostgresRepository)

	db, err := sql.Open("pgx", cfg.DBUrl)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(maxOpenConns)
	db.SetMaxIdleConns(maxIdleConns)
	db.SetConnMaxLifetime(connMaxLifetime)
	db.SetConnMaxIdleTime(connMaxIdleTime)

	repo.QBQueries = New(db)
	cfg.AddCloser(db.Close)

	return repo, nil
}

func (r *PostgresRepository) Create(ctx context.Context, url string) (*domain.ShortUrl, error) {
	if res, err := domain.ValidateURL(url); !res {
		return nil, fmt.Errorf("invalid URL: %w", err)
	}

	code, err := domain.GenerateCode(url)
	if err != nil {
		return nil, fmt.Errorf("error generating code: %w", err)
	}

	entry, err := r.QBQueries.CreateShortUrl(ctx, CreateShortUrlParams{
		Code:        code,
		OriginalUrl: url,
	})
	if err != nil {
		return nil, err
	}

	return entryToDomain(entry), nil
}

func (r *PostgresRepository) FindByCode(ctx context.Context, code string) (*domain.ShortUrl, error) {
	entry, err := r.QBQueries.GetByCode(ctx, code)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}

	return entryToDomain(entry), nil
}

func (r *PostgresRepository) Delete(ctx context.Context, code string) error {
	return r.QBQueries.DeleteByCode(ctx, code)
}

func (r *PostgresRepository) List(ctx context.Context, page int, amount int, search string) ([]*domain.ShortUrl, int, error) {
	offset := amount * (page - 1)
	totalAmount, err := r.QBQueries.Amount(ctx)
	if err != nil {
		return nil, 0, err
	}

	var result []ShortUrl
	if search == "" {
		result, err = r.QBQueries.List(ctx, ListParams{
			Offset: int32(offset),
			Limit:  int32(amount),
		})
	} else {
		result, err = r.QBQueries.Search(ctx, SearchParams{
			Offset:      int32(offset),
			Limit:       int32(amount),
			OriginalUrl: search,
		})
	}
	if err != nil {
		return nil, 0, err
	}

	var retVal []*domain.ShortUrl
	for _, val := range result {
		retVal = append(retVal, entryToDomain(val))
	}

	return retVal, int(totalAmount), nil
}

func (r *PostgresRepository) IncrementClicks(ctx context.Context, code string) (*domain.ShortUrl, error) {
	entry, err := r.QBQueries.IncrementClicks(ctx, code)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}

	return entryToDomain(entry), nil
}
