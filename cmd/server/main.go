package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

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

	api := NewApiController(srv, cfg.Domain)

	router := gin.Default()

	router.Use(ErrorHandler())

	router.GET("/healthz", api.GetHealth)
	router.GET("/:code", api.Redirect)

	apiRoute := router.Group("/api")
	{
		v1Route := apiRoute.Group("/v1")
		{
			linksRoute := v1Route.Group("/links")
			{
				linksRoute.GET("/:code/stats", api.GetStats)
				linksRoute.GET("/:code", api.GetSingle)
				linksRoute.POST("", api.Create)
				linksRoute.DELETE("/:code", api.Remove)
			}
		}
	}

	server := http.Server{
		Addr:    cfg.Port,
		Handler: router,
	}

	go func() {
		log.Printf("Starting server on %s\n", cfg.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Listen error: %s", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutdown signal recieved")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Forced shutdown: %s", err)
	}

	log.Println("Server exited gracefully")
}
