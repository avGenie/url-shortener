package config

import (
	"flag"
	"fmt"

	"github.com/caarlos0/env/v10"
)

var (
	Config config
)

type config struct {
	NetAddr       string `env:"SERVER_ADDRESS"`
	BaseURIPrefix string `env:"BASE_URL"`
}

func init() {
	flag.StringVar(&Config.NetAddr, "a", "localhost:8080", "net address host:port")
	flag.StringVar(&Config.BaseURIPrefix, "b", "http://localhost:8080", "base output short URL")
}

// Parses config from environment and command line arguments
func ParseConfig() {
	flag.Parse()

	if err := env.Parse(&Config); err != nil {
		fmt.Printf("%+v\n", err)
	}

	fmt.Printf("%+v\n", Config)
}
