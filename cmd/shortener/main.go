package main

import (
	"net/http"

	"github.com/avGenie/url-shortener/internal/app/config"
	"github.com/avGenie/url-shortener/internal/app/handlers"
)

func main() {
	config.ParseConfig()

	err := http.ListenAndServe(config.Config.NetAddr, handlers.CreateRouter())
	if err != nil && err != http.ErrServerClosed {
		panic(err.Error())
	}
}
