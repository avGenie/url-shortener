package main

import (
	"flag"
	"net/http"

	"github.com/avGenie/url-shortener/internal/app/config"
	"github.com/avGenie/url-shortener/internal/app/handlers"
)

func main() {
	flag.Parse()
	err := http.ListenAndServe(config.NetAddr, handlers.CreateRouter())
	if err != nil {
		panic(err.Error())
	}
}
