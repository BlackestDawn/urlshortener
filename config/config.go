package config

import (
	"log"
	"os"
)

type Config struct {
	DBUrl string
	Port  string
	Env   string
}

func New() Config {
	dbUrl := os.Getenv("DATABASE_URL")
	if dbUrl == "" {
		log.Fatalln("Missing URL for database connection (DATABASE_URL)")
	}

	listenPort := os.Getenv("LISTEN_PORT")
	if listenPort == "" {
		listenPort = defaultListenPort
	}

	appEnv := os.Getenv("PLATFORM_ENV")
	if appEnv == "" {
		appEnv = defaultAppEnv
	}

	return Config{
		DBUrl: dbUrl,
		Port:  listenPort,
		Env:   appEnv,
	}
}
