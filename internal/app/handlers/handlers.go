package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/avGenie/url-shortener/internal/app/config"
	"github.com/avGenie/url-shortener/internal/app/entity"
	"github.com/avGenie/url-shortener/internal/app/models"
	"github.com/avGenie/url-shortener/internal/app/storage"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type favContextKey string

const (
	maxEncodedSize = 8

	baseURIPrefixCtx = favContextKey("baseURIPrefix")
)

var (
	urls *storage.URLStorage
)

func InitStorage(config config.Config) error {
	storage, err := storage.NewURLStorage(config.DBFileStoragePath)
	urls = storage

	return err
}

func CloseStorage(config config.Config) {
	if strings.Contains(config.DBFileStoragePath, os.TempDir()) {
		os.Remove(config.DBFileStoragePath)
	}
}

// Processes POST request. Sends short URL in http://localhost:8080/id format.
//
// Encodes given URL using base64 encoding scheme and puts it to the URL's map.
//
// Returns 201 status code if processing was successfull, otherwise returns 400.
func PostHandlerURL(baseURIPrefix string, writer http.ResponseWriter, req *http.Request) {
	zap.L().Debug("POST handler URL processing")

	inputURL, err := io.ReadAll(req.Body)
	defer req.Body.Close()

	if err != nil {
		zap.L().Error(CannotProcessURL, zap.Error(err))
		http.Error(writer, CannotProcessURL, http.StatusBadRequest)
		return
	}

	if ok := entity.IsValidURL(string(inputURL)); !ok {
		zap.L().Error(WrongURLFormat, zap.Error(err))
		http.Error(writer, WrongURLFormat, http.StatusBadRequest)
		return
	}

	if baseURIPrefix == "" {
		zap.L().Error("invalid base URI prefix", zap.String("base URI prefix", baseURIPrefix))
		http.Error(writer, InternalServerError, http.StatusInternalServerError)
		return
	}

	outputURL, err := postURLProcessing(string(inputURL), baseURIPrefix)
	if err != nil || outputURL == "" {
		zap.L().Error("could not create a short URL", zap.String("error", err.Error()))
		http.Error(writer, InternalServerError, http.StatusInternalServerError)
		return
	}

	zap.L().Info("url has been created succeessfully", zap.String("output url", outputURL))

	writer.Header().Set("Content-Type", "text/plain; charset=utf-8")
	writer.WriteHeader(http.StatusCreated)
	io.WriteString(writer, outputURL)
}

func PostHandlerJSON(baseURIPrefix string, writer http.ResponseWriter, req *http.Request) {
	zap.L().Debug("POST handler JSON processing")

	inputRequest := &models.Request{}
	err := json.NewDecoder(req.Body).Decode(&inputRequest)
	defer req.Body.Close()
	if err != nil {
		zap.L().Error(CannotProcessJSON, zap.Error(err))
		http.Error(writer, CannotProcessJSON, http.StatusBadRequest)
		return
	}

	if ok := entity.IsValidURL(inputRequest.URL); !ok {
		zap.L().Error(WrongURLFormat, zap.Error(err))
		http.Error(writer, WrongURLFormat, http.StatusBadRequest)
		return
	}

	if baseURIPrefix == "" {
		zap.L().Error("invalid base URI prefix", zap.String("base URI prefix", baseURIPrefix))
		http.Error(writer, InternalServerError, http.StatusInternalServerError)
		return
	}

	outputURL, err := postURLProcessing(inputRequest.URL, baseURIPrefix)
	if err != nil || outputURL == "" {
		zap.L().Error("could not create a short URL", zap.String("error", err.Error()))
		http.Error(writer, InternalServerError, http.StatusInternalServerError)
		return
	}

	response := &models.Response{
		URL: outputURL,
	}
	zap.L().Info("url has been created succeessfully", zap.String("output url", response.URL))

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusCreated)
	if err = json.NewEncoder(writer).Encode(response); err != nil {
		zap.L().Error("invalid response", zap.Any("response", response))
		http.Error(writer, "internal server error", http.StatusInternalServerError)
		return
	}

	zap.L().Debug("sending HTTP 200 response")
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
		zap.L().Error(ShortURLNotInDB)
		http.Error(writer, ShortURLNotInDB, http.StatusBadRequest)
		return
	}

	zap.L().Info("url has been decoded succeessfully", zap.String("decoded url", decodedURL.String()))

	writer.Header().Set("Content-Type", "text/plain; charset=utf-8")
	writer.Header().Set("Location", decodedURL.String())
	writer.WriteHeader(http.StatusTemporaryRedirect)
}
