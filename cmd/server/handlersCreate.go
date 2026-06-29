package main

import (
	"net/http"

	"github.com/BlackestDawn/urlshortener/api"
	"github.com/gin-gonic/gin"
)

func (a *ApiController) Create(c *gin.Context) {
	var data api.UrlDto
	err := c.ShouldBindJSON(&data)
	if err != nil {
		c.Error(err)
		return
	}

	code, err := a.srv.Shorten(data.Url)
	if err != nil {
		c.Error(err)
		return
	}

	url := "https://" + a.domain + "/" + code
	c.JSON(http.StatusCreated, api.ShortenedUrlDto{ShortenedUrl: url})
}
