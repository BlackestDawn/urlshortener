package api

import "time"

type UrlDto struct {
	Url string `json:"url"`
}

type ShortenedUrlDto struct {
	ShortenedUrl string `json:"shortenedUrl"`
}

type UrlStatsDto struct {
	Url       string    `json:"url"`
	Code      string    `json:"code"`
	Hits      int       `json:"hits"`
	CreatedAt time.Time `json:"createdAt"`
}
