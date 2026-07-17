package main

import (
	"errors"
	"net/http"

	"github.com/BlackestDawn/urlshortener/api"
	"github.com/BlackestDawn/urlshortener/internal/domain"
	"github.com/gin-gonic/gin"
)

func (a *ApiController) Create(c *gin.Context) {
	var data api.UrlDto
	err := c.ShouldBindJSON(&data)
	if err != nil {
		var maxBytesErr *http.MaxBytesError
		if errors.As(err, &maxBytesErr) {
			c.Error(domain.ErrRequestTooLarge)
			return
		}
		c.Error(domain.ErrInvalidJson)
		return
	}

	code, err := a.srv.Shorten(c.Request.Context(), data.Url)
	if err != nil {
		c.Error(err)
		return
	}

	url := "https://" + a.domain + "/" + code
	c.JSON(http.StatusCreated, api.ShortenedUrlDto{ShortenedUrl: url})
}
