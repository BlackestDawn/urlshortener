package main

import (
	"net/http"

	"github.com/BlackestDawn/urlshortener/api"
	"github.com/gin-gonic/gin"
)

func (a *ApiController) GetSingle(c *gin.Context) {
	code := c.Param("code")
	url, err := a.srv.Resolve(code)
	if err != nil {
		c.Error(err)
	}

	c.JSON(http.StatusOK, api.UrlDto{Url: url})
}

func (a *ApiController) GetStats(c *gin.Context) {
	code := c.Param("code")
	entity, err := a.srv.GetStats(code)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, api.EntityToStatDto(entity))
}
