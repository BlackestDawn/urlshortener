package main

import (
	"fmt"
	"net/http"

	"github.com/BlackestDawn/urlshortener/internal/domain"
	"github.com/gin-gonic/gin"
)

func (a *ApiController) Redirect(c *gin.Context) {
	code := c.Param("code")
	url, err := a.srv.Resolve(code)
	if err != nil {
		if err == domain.ErrNotFound {
			respondJSONError(c, http.StatusNotFound, fmt.Sprintf("url for code '%s' not found", code), err)
		}
		respondJSONError(c, http.StatusBadRequest, "could not complete request", err)
	}

	c.Redirect(http.StatusPermanentRedirect, url)
}
