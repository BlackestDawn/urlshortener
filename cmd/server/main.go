package main

import (
	"log"

	"github.com/BlackestDawn/urlshortener/config"
	"github.com/BlackestDawn/urlshortener/internal/repository"
	"github.com/BlackestDawn/urlshortener/internal/service"
	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.NewConfig()
	defer cfg.Cleanup()

	repo, err := repository.NewPGRepository(cfg)
	if err != nil {
		log.Fatalf("failed setting up repo: %s", err.Error())
	}

	srv := service.NewShortenService(repo)

	api := NewApiController(srv)

	router := gin.Default()

	router.GET("/healthz", api.GetHealth)

	router.Run(cfg.Port)
}
