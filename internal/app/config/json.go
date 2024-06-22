package config

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

// jsonConfig JSON config struct
type jsonConfig struct {
	NetAddr           string `json:"server_address"`
	BaseURIPrefix     string `json:"base_url"`
	DBFileStoragePath string `json:"file_storage_path"`
	DBStorageConnect  string `json:"database_dsn"`
	Enable_HTTPS      bool   `json:"enable_https"`
}

func parseJSONConfig(config *Config) error {
	if config.ConfigFile == "" {
		return fmt.Errorf("JSON file config not found")
	}

	jsonFile, err := os.Open(config.ConfigFile)
	if err != nil {
		return fmt.Errorf("couldn't open JSON config file %w", err)
	}

	defer jsonFile.Close()

	byteValue, err := io.ReadAll(jsonFile)
	if err != nil {
		return fmt.Errorf("couldn't read data from JSON config file %w", err)
	}

	var jsonConf jsonConfig
	err = json.Unmarshal(byteValue, &jsonConf)
	if err != nil {
		return fmt.Errorf("couldn't read data from JSON config file %w", err)
	}

	fillConfigByJSON(config, jsonConf)

	return nil
}

func fillConfigByJSON(config *Config, jsonConfig jsonConfig) {
	if config.NetAddr == defaultNetAddr {
		config.NetAddr = jsonConfig.NetAddr
	}

	if config.BaseURIPrefix == defaultBaseURIPrefix {
		config.BaseURIPrefix = jsonConfig.BaseURIPrefix
	}

	if config.DBFileStoragePath == defaultFileStoragePath {
		config.DBFileStoragePath = jsonConfig.DBFileStoragePath
	}

	if config.DBStorageConnect == "" {
		config.DBStorageConnect = jsonConfig.DBStorageConnect
	}

	if !config.Enable_HTTPS {
		config.Enable_HTTPS = jsonConfig.Enable_HTTPS
	}
}
