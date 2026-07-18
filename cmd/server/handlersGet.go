package main

import (
	"net/http"

	"github.com/BlackestDawn/urlshortener/api"
	"github.com/gin-gonic/gin"
)

func (a *ApiController) GetSingle(c *gin.Context) {
	code := c.Param("code")
	entry, err := a.srv.GetStats(c.Request.Context(), code)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, api.UrlDto{Url: entry.OriginalUrl})
}

func (a *ApiController) GetStats(c *gin.Context) {
	code := c.Param("code")
	entity, err := a.srv.GetStats(c.Request.Context(), code)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, api.EntityToStatDto(entity))
}
