package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (a *ApiController) Remove(c *gin.Context) {
	code := c.Param("code")
	err := a.srv.Delete(code)
	if err != nil {
		respondJSONError(c, http.StatusInternalServerError, "something went wrong removing entry", err)
	}

	respondJSON(c, http.StatusNoContent, nil)
}
