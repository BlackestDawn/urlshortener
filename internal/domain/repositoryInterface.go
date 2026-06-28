package domain

type IRepository interface {
	Create(url string) (*ShortUrl, error)
	FindByCode(code string) (*ShortUrl, error)
	IncrementClicks(code string) error
	List(page int, amount int, search string) ([]*ShortUrl, int, error)
	Delete(code string) error
}
