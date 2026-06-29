package main

import (
	"errors"
	"net/http"

	"github.com/BlackestDawn/urlshortener/internal/domain"
	"github.com/gin-gonic/gin"
)

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		// process requests first
		c.Next()

		if len(c.Errors) > 0 {
			err := c.Errors.Last().Err

			type errVal struct {
				Error string
			}

			if errors.Is(err, domain.ErrInvalidJson) {
				c.JSON(http.StatusUnprocessableEntity, errVal{Error: err.Error()})
				return
			}

			if errors.Is(err, domain.ErrInvalidUrl) {
				c.JSON(http.StatusBadRequest, errVal{Error: err.Error()})
				return
			}

			if errors.Is(err, domain.ErrNotFound) {
				c.JSON(http.StatusNotFound, errVal{Error: err.Error()})
				return
			}

			c.JSON(http.StatusInternalServerError, errVal{Error: "Internal server error"})
		}
	}
}
