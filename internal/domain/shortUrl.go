package domain

import (
	"time"

	"github.com/google/uuid"
)

type ShortUrl struct {
	ID          uuid.UUID
	Code        string
	OriginalUrl string
	CreatedAt   time.Time
	Clicks      int
}
