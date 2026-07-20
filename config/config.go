package config

import (
	"log"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

type Config struct {
	DBUrl   string
	Port    string
	Env     string
	Domain  string
	closers []func() error
}

func NewConfig() *Config {
	err := godotenv.Load(findEnvFile(""))
	if err != nil {
		slog.Info("error loading env file", "error", err.Error())
	}

	appEnv := os.Getenv("URLSHORTENER_ENV")
	if appEnv == "" {
		appEnv = defaultAppEnv
	}

	err = godotenv.Load(findEnvFile(appEnv))
	if err != nil {
		slog.Info("error loading env file", "error", err.Error())
	}

	dbUrl := os.Getenv("DATABASE_URL")
	if dbUrl == "" {
		log.Fatalln("Missing URL for database connection (DATABASE_URL)")
	}

	listenPort := os.Getenv("LISTEN_PORT")
	if listenPort == "" {
		listenPort = defaultListenPort
	}

	domain := os.Getenv("SHORTEN_DOMAIN")
	if domain == "" {
		domain = defaultDomain
	}

	return &Config{
		DBUrl:  dbUrl,
		Port:   ":" + listenPort,
		Env:    appEnv,
		Domain: domain,
	}
}

func (c *Config) AddCloser(closer func() error) {
	c.closers = append(c.closers, closer)
}

func (c *Config) Cleanup() {
	logger := slog.Default()
	for _, closer := range c.closers {
		err := closer()
		if err != nil {
			logger.Warn("Something went wrong during shutdown", "error", err.Error())
		}
	}
}

func findEnvFile(fileNameExtension string) string {
	dir, _ := os.Getwd()

	for {
		path := filepath.Join(dir, ".env."+fileNameExtension)
		if _, err := os.Stat(path); err == nil {
			return path
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}

	return ""
}
