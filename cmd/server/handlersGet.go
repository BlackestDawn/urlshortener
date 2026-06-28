package main

import (
	"fmt"
	"net/http"

	"github.com/BlackestDawn/urlshortener/internal/domain"
	"github.com/gin-gonic/gin"
)

func (a *ApiController) GetSingle(c *gin.Context) {
	code := c.Param("code")
	url, err := a.srv.Resolve(code)
	if err != nil {
		if err == domain.ErrNotFound {
			respondJSONError(c, http.StatusNotFound, fmt.Sprintf("url for code '%s' not found", code), err)
		}
		respondJSONError(c, http.StatusBadRequest, "could not complete request", err)
	}

	type response struct {
		Url string
	}

	respondJSON(c, http.StatusOK, response{Url: url})
}

func (a *ApiController) GetStats(c *gin.Context) {
	code := c.Param("code")
	entity, err := a.srv.GetStats(code)
	if err != nil {
		if err == domain.ErrNotFound {
			respondJSONError(c, http.StatusNotFound, fmt.Sprintf("url for code '%s' not found", code), err)
		}
		respondJSONError(c, http.StatusBadRequest, "could not complete request", err)
	}

	respondJSON(c, http.StatusOK, entity)
}
