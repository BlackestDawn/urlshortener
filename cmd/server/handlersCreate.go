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
	err := c.ShouldBindJSON(data)
	if err != nil {
		respondJSONError(c, http.StatusInternalServerError, "error decoding data", err)
	}

	code, err := a.srv.Shorten(data.Url)
	if err != nil {
		if errors.Is(err, domain.ErrInvalidUrl) {
			respondJSONError(c, http.StatusBadRequest, "invalid URL", err)
			return
		}
		respondJSONError(c, http.StatusInternalServerError, "error shortening URL", err)
	}

	url := "https://" + a.domain + "/" + code
	respondJSON(c, http.StatusOK, api.ShortenedUrlDto{ShortenedUrl: url})
}
