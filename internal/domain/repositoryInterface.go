package domain

type IRepository interface {
	Create(entry ShortUrl) error
	FindByCode(code string) ShortUrl
	IncrementClicks(code string) (int, error)
	List(skip int, take int, searchParam string) (*[]ShortUrl, error)
	Delete(code string) error
}
