package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (a *ApiController) Redirect(c *gin.Context) {
	code := c.Param("code")
	url, err := a.srv.Resolve(code)
	if err != nil {
		c.Error(err)
		return
	}

	c.Redirect(http.StatusPermanentRedirect, url)
}
