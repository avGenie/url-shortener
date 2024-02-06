package handlers

import (
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/avGenie/url-shortener/internal/app/config"
	"github.com/go-chi/chi/v5"
)

const (
	maxEncodedSize = 8
)

var (
	urls = make(map[string]string)

	ErrInvalidGivenURL = errors.New("given URL is mapped")

	EmptyURL        = "URL is empty"
	ShortURLNotInDB = "given short URL did not find in database"
)

// Processes POST request. Sends short URL in http://localhost:8080/id format.
//
// Encodes given URL using base64 encoding scheme and puts it to the URL's map.
//
// Returns 201 status code if processing was successfull, otherwise returns 400.
func PostHandler(writer http.ResponseWriter, req *http.Request) {
	fmt.Println("PostHandler")
	writer.Header().Set("Content-Type", "text/plain")
	bodyBytes, err := io.ReadAll(req.Body)
	if err != nil {
		fmt.Println("read error")
		IncorrectRequestHandler(writer, req, fmt.Sprintf("cannot process URL: %s", err.Error()))
		return
	}

	if len(bodyBytes) == 0 {
		fmt.Println("len error")
		IncorrectRequestHandler(writer, req, EmptyURL)
		return
	}

	bodyString := string(bodyBytes)
	encodedURL := base64.StdEncoding.EncodeToString(bodyBytes)
	shortURL := encodedURL[:maxEncodedSize]

	urls[shortURL] = bodyString

	outputURL := fmt.Sprintf("%s/%s", config.BaseURIPrefix, shortURL)
	fmt.Println(outputURL)

	writer.WriteHeader(http.StatusCreated)
	io.WriteString(writer, outputURL)
}

// Processes GET request. Sends the source address at the given short address
//
// # Sends short URL back to the original using from the URL's map
//
// Returns 307 status code if processing was successfull, otherwise returns 400.
func GetHandler(writer http.ResponseWriter, req *http.Request) {
	writer.Header().Set("Content-Type", "text/plain")

	shortURL := chi.URLParam(req, "url")
	fmt.Printf("shortURL: %s\n", shortURL)

	decodedURL, ok := urls[shortURL]
	if !ok {
		fmt.Println("GetHandler error")
		IncorrectRequestHandler(writer, req, ShortURLNotInDB)
		return
	}

	fmt.Printf("decodedURL: %s\n", decodedURL)

	writer.Header().Set("Location", decodedURL)
	writer.WriteHeader(http.StatusTemporaryRedirect)
}

// Sends 400 status code with error message
func IncorrectRequestHandler(writer http.ResponseWriter, req *http.Request, message string) {
	writer.WriteHeader(http.StatusBadRequest)
	io.WriteString(writer, message)
}
