package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/BlackestDawn/urlshortener/config"
	"github.com/BlackestDawn/urlshortener/internal/domain"
	_ "github.com/jackc/pgx/v5/stdlib"
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

	repo.QBQueries = New(db)
	cfg.AddCloser(db.Close)

	return repo, nil
}

func (r *PostgresRepository) Create(url string) (*domain.ShortUrl, error) {
	if res, err := domain.ValidateURL(url); !res {
		return nil, fmt.Errorf("invalid URL: %w", err)
	}

	code, err := domain.GenerateCode(url)
	if err != nil {
		return nil, fmt.Errorf("error generating code: %w", err)
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

func (r *PostgresRepository) FindByCode(code string) (*domain.ShortUrl, error) {
	entry, err := r.QBQueries.GetByCode(context.Background(), code)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}

	return entryToDomain(entry), nil
}

func (r *PostgresRepository) Delete(code string) error {
	return r.QBQueries.DeleteByCode(context.Background(), code)
}

func (r *PostgresRepository) List(page int, amount int, search string) ([]*domain.ShortUrl, int, error) {
	offset := amount * (page - 1)
	totalAmount, err := r.QBQueries.Amount(context.Background())
	if err != nil {
		return nil, 0, err
	}

	var result []ShortUrl
	if search == "" {
		result, err = r.QBQueries.List(context.Background(), ListParams{
			Offset: int32(offset),
			Limit:  int32(amount),
		})
	} else {
		result, err = r.QBQueries.Search(context.Background(), SearchParams{
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

func (r *PostgresRepository) IncrementClicks(code string) error {
	res, err := r.QBQueries.GetByCode(context.Background(), code)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.ErrNotFound
		}
		return err
	}

	err = r.QBQueries.IncrementClicks(context.Background(), IncrementClicksParams{
		Code:   code,
		Clicks: res.Clicks + 1,
	})

	return err
}
