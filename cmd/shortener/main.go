package main

import (
	"net/http"

	"github.com/avGenie/url-shortener/internal/app/handlers"
)

func main() {
	err := http.ListenAndServe(":8080", handlers.CreateRouter())
	if err != nil {
		panic(err.Error())
	}
}
