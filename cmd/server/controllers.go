package main

import (
	"net/http"

	"github.com/BlackestDawn/urlshortener/internal/service"
	"github.com/gin-gonic/gin"
)

type ApiController struct {
	srv service.IShorten
}

func NewApiController(srv service.IShorten) *ApiController {
	return &ApiController{
		srv: srv,
	}
}

func (a *ApiController) GetHealth(c *gin.Context) {
	c.String(http.StatusOK, "you hit endpoint: get health")
}
