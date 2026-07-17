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
	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
	ratelimiter "github.com/rleungx/gin-ratelimiter"
)

const (
	maxRequestBodyBytes = 4 << 10 // 4 KiB, generous for a single URL payload
	requestTimeout      = 10 * time.Second
)

func main() {
	cfg := config.NewConfig()

	repo, err := repository.NewPGRepository(cfg)
	if err != nil {
		log.Fatalf("failed setting up repo: %s", err.Error())
	}

	srv := service.NewShortenService(repo)

	api := NewApiController(srv, cfg.Domain)
	limiter := ratelimiter.New()
	limiter.UpdateRateLimit("/healthz", 3, 3)
	limiter.UpdateConcurrencyLimit("/healthz", 3)

	router := gin.New()

	router.Use(requestid.New())
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(MaxBodySize(maxRequestBodyBytes))
	router.Use(limiter.Middleware(ratelimiter.WithRateLimit(30, 60), ratelimiter.WithConcurrencyLimit(5)))
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
		Handler: http.TimeoutHandler(router, requestTimeout, "request timed out"),
	}

	serverErr := make(chan error, 1)
	go func() {
		log.Printf("Starting server on %s\n", cfg.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			serverErr <- err
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-serverErr:
		log.Printf("Listen error: %s", err)
		cfg.Cleanup()
		os.Exit(1)
	case <-quit:
		log.Println("Shutdown signal recieved")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Printf("Forced shutdown: %s", err)
		cfg.Cleanup()
		os.Exit(1)
	}

	cfg.Cleanup()
	log.Println("Server exited gracefully")
}
