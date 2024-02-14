package handlers

import (
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/avGenie/url-shortener/internal/app/entity"
	"github.com/avGenie/url-shortener/internal/app/storage"
	"github.com/go-chi/chi/v5"
)

const (
	maxEncodedSize = 8
)

var (
	urls = storage.NewURLStorage()

	EmptyURL         = "URL is empty"
	WrongURLFormat   = "wrong URL format"
	ShortURLNotInDB  = "given short URL did not find in database"
	CannotProcessURL = "cannot process URL"
)

// Processes POST request. Sends short URL in http://localhost:8080/id format.
//
// Encodes given URL using base64 encoding scheme and puts it to the URL's map.
//
// Returns 201 status code if processing was successfull, otherwise returns 400.
func PostHandler(baseURIPrefix string, writer http.ResponseWriter, req *http.Request) {
	userURL, err := io.ReadAll(req.Body)
	defer req.Body.Close()

	if err != nil {
		http.Error(writer, CannotProcessURL, http.StatusBadRequest)
		return
	}

	if len(userURL) == 0 {
		http.Error(writer, EmptyURL, http.StatusBadRequest)
		return
	}

	if ok := entity.IsValidURL(string(userURL)); !ok {
		http.Error(writer, WrongURLFormat, http.StatusBadRequest)
		return
	}

	encodedURL := base64.StdEncoding.EncodeToString(userURL)
	shortURL := encodedURL[:maxEncodedSize]

	urls.Add(*entity.ParseURL(shortURL), *entity.ParseURL(string(userURL)))

	outputURL := fmt.Sprintf("%s/%s", baseURIPrefix, shortURL)

	log.Println("Created URL: ", outputURL)

	writer.Header().Set("Content-Type", "text/plain; charset=utf-8")
	writer.WriteHeader(http.StatusCreated)
	io.WriteString(writer, outputURL)
}

// Processes GET request. Sends the source address at the given short address
//
// # Sends short URL back to the original using from the URL's map
//
// Returns 307 status code if processing was successfull, otherwise returns 400.
func GetHandler(writer http.ResponseWriter, req *http.Request) {
	shortURL := chi.URLParam(req, "url")
	decodedURL, ok := urls.Get(*entity.ParseURL(shortURL))

	if !ok {
		http.Error(writer, ShortURLNotInDB, http.StatusBadRequest)
		return
	}

	log.Println("Decoded URL: ", decodedURL)

	writer.Header().Set("Content-Type", "text/plain; charset=utf-8")
	writer.Header().Set("Location", decodedURL.String())
	writer.WriteHeader(http.StatusTemporaryRedirect)
}
