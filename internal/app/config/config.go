package config

import (
	"flag"

	"github.com/caarlos0/env/v10"
)

type Config struct {
	NetAddr       string `env:"SERVER_ADDRESS"`
	BaseURIPrefix string `env:"BASE_URL"`
}

func InitConfig() (config Config) {
	flag.StringVar(&config.NetAddr, "a", "localhost:8080", "net address host:port")
	flag.StringVar(&config.BaseURIPrefix, "b", "http://localhost:8080", "base output short URL")
	flag.Parse()

	if err := env.Parse(&config); err != nil {
		panic(err.Error())
	}

	return
}
