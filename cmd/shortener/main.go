package main

import (
	"net/http"

	handlers "github.com/avGenie/url-shortener/internal/app/handlers"
)

func main() {
	err := run()
	if err != nil {
		panic(err.Error())
	}
}

// Runs HTTP-Server
func run() error {
	mux := createMux()
	return http.ListenAndServe(":8080", mux)
}

// Creates new ServeMux
func createMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", webhook)

	return mux
}

// HTTP POST and GET requests handler. If there is another request, returns 400 status code
func webhook(writer http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodPost:
		handlers.PostHandler(writer, req)
	case http.MethodGet:
		handlers.GetHandler(writer, req)
	default:
		handlers.IncorrectRequestHandler(writer, req, "incorrect request")
	}
}
