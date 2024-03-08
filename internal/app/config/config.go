package config

import (
	"flag"

	"github.com/caarlos0/env/v10"
)

type Config struct {
	NetAddr       string `env:"SERVER_ADDRESS"`
	BaseURIPrefix string `env:"BASE_URL"`
	LogLevel      string `env:"LOG_LEVEL"`
	DBStorage     string `env:"FILE_STORAGE_PATH"`
}

func InitConfig() (config Config) {
	flag.StringVar(&config.NetAddr, "a", "localhost:8080", "net address host:port")
	flag.StringVar(&config.BaseURIPrefix, "b", "http://localhost:8080", "base output short URL")
	flag.StringVar(&config.LogLevel, "l", "info", "log level")
	flag.StringVar(&config.DBStorage, "f", "/tmp/short-url-db.json", "database storage path")
	flag.Parse()

	if err := env.Parse(&config); err != nil {
		panic(err.Error())
	}

	return
}
