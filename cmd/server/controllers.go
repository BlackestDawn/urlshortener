package main

import (
	"net/http"

	"github.com/BlackestDawn/urlshortener/internal/service"
	"github.com/gin-gonic/gin"
)

type ApiController struct {
	srv    service.IShorten
	domain string
}

func NewApiController(srv service.IShorten, domain string) *ApiController {
	return &ApiController{
		srv:    srv,
		domain: domain,
	}
}

func (a *ApiController) GetHealth(c *gin.Context) {
	type response struct {
		Health string
	}
	c.JSON(http.StatusOK, response{Health: "OK"})
}
