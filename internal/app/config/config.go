// Package config implements application config
package config

import (
	"flag"
	"fmt"

	"github.com/caarlos0/env/v10"
)

const (
	defaultNetAddr         = "localhost:8080"
	defaultBaseURIPrefix   = "http://localhost:8080"
	defaultLogLevel        = "debug"
	defaultFileStoragePath = "/tmp/short-url-db.json"
)

// Config struct
type Config struct {
	NetAddr           string `json:"server_address" env:"SERVER_ADDRESS"`
	BaseURIPrefix     string `json:"base_url" env:"BASE_URL"`
	LogLevel          string `json:"-" env:"LOG_LEVEL"`
	DBFileStoragePath string `json:"file_storage_path" env:"FILE_STORAGE_PATH"`
	DBStorageConnect  string `json:"database_dsn" env:"DATABASE_DSN"`
	ProfilerFile      string `json:"-" env:"PROFILER_FILE"`
	ConfigFile        string `json:"-" env:"CONFIG"`
	EnableHTTPS       bool   `json:"enable_https" env:"ENABLE_HTTPS"`
}

// InitConfig Initialize config from flag and env variables
func InitConfig() (Config, error) {
	var config Config
	flag.StringVar(&config.NetAddr, "a", defaultNetAddr, "net address host:port")
	flag.StringVar(&config.BaseURIPrefix, "b", defaultBaseURIPrefix, "base output short URL")
	flag.StringVar(&config.LogLevel, "l", defaultLogLevel, "log level")
	flag.StringVar(&config.DBFileStoragePath, "f", defaultFileStoragePath, "database storage path")
	flag.StringVar(&config.DBStorageConnect, "d", "", "database credentials in format: host=host port=port user=myuser password=xxxx dbname=mydb sslmode=disable")
	flag.StringVar(&config.ProfilerFile, "p", "", "profiler file name")
	flag.StringVar(&config.ConfigFile, "c", "", "configuration JSON file")
	flag.BoolVar(&config.EnableHTTPS, "s", false, "enable HTTPS")
	flag.Parse()

	if err := env.Parse(&config); err != nil {
		return Config{}, err
	}

	if err := parseJSONConfig(&config); err != nil {
		fmt.Printf("couldn't parse config file: %s\n", err.Error())
	}

	return config, nil
}
