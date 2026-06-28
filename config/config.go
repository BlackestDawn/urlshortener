package config

import (
	"log"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

type Config struct {
	DBUrl   string
	Port    string
	Env     string
	closers []func() error
}

func NewConfig() *Config {
	godotenv.Load(findEnvFile(""))

	appEnv := os.Getenv("URLSHOTRENER_ENV")
	if appEnv == "" {
		appEnv = defaultAppEnv
	}

	godotenv.Load(findEnvFile(appEnv))

	dbUrl := os.Getenv("DATABASE_URL")
	if dbUrl == "" {
		log.Fatalln("Missing URL for database connection (DATABASE_URL)")
	}

	listenPort := os.Getenv("LISTEN_PORT")
	if listenPort == "" {
		listenPort = defaultListenPort
	}

	return &Config{
		DBUrl: dbUrl,
		Port:  ":" + listenPort,
		Env:   appEnv,
	}
}

func (c *Config) AddCloser(closer func() error) {
	c.closers = append(c.closers, closer)
}

func (c *Config) Cleanup() {
	for _, closer := range c.closers {
		closer()
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
