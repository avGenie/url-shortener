package handlers

import (
	"encoding/base64"
	"fmt"
	"io"
	"net/http"

	"github.com/avGenie/url-shortener/internal/app/entity"
	"github.com/avGenie/url-shortener/internal/app/logger"
	"github.com/avGenie/url-shortener/internal/app/storage"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
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
		logger.Log.Error(CannotProcessURL, zap.Error(err))
		http.Error(writer, CannotProcessURL, http.StatusBadRequest)
		return
	}

	if len(userURL) == 0 {
		logger.Log.Error(EmptyURL, zap.Error(err))
		http.Error(writer, EmptyURL, http.StatusBadRequest)
		return
	}

	if ok := entity.IsValidURL(string(userURL)); !ok {
		logger.Log.Error(WrongURLFormat, zap.Error(err))
		http.Error(writer, WrongURLFormat, http.StatusBadRequest)
		return
	}

	encodedURL := base64.StdEncoding.EncodeToString(userURL)
	shortURL := encodedURL[:maxEncodedSize]

	urls.Add(*entity.ParseURL(shortURL), *entity.ParseURL(string(userURL)))

	outputURL := fmt.Sprintf("%s/%s", baseURIPrefix, shortURL)

	logger.Log.Info("url has been created succeessfully", zap.String("output url", outputURL))

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
		logger.Log.Error(ShortURLNotInDB)
		http.Error(writer, ShortURLNotInDB, http.StatusBadRequest)
		return
	}

	logger.Log.Info("url has been decoded succeessfully", zap.String("decoded url", decodedURL.String()))

	writer.Header().Set("Content-Type", "text/plain; charset=utf-8")
	writer.Header().Set("Location", decodedURL.String())
	writer.WriteHeader(http.StatusTemporaryRedirect)
}
