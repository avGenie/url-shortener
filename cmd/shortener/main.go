package main

import (
	"net/http"

	"github.com/avGenie/url-shortener/internal/app/config"
	"github.com/avGenie/url-shortener/internal/app/handlers"
)

func main() {
	cnf := config.InitConfig()

	err := http.ListenAndServe(cnf.NetAddr, handlers.CreateRouter(cnf))
	if err != nil && err != http.ErrServerClosed {
		panic(err.Error())
	}
}
