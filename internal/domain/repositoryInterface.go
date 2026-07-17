package domain

import "context"

type IRepository interface {
	Create(ctx context.Context, url string) (*ShortUrl, error)
	FindByCode(ctx context.Context, code string) (*ShortUrl, error)
	IncrementClicks(ctx context.Context, code string) error
	List(ctx context.Context, page int, amount int, search string) ([]*ShortUrl, int, error)
	Delete(ctx context.Context, code string) error
}
