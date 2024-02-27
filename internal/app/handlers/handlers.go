package handlers

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/avGenie/url-shortener/internal/app/entity"
	"github.com/avGenie/url-shortener/internal/app/logger"
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
	urls = storage.NewURLStorage()
)

// Processes POST request. Sends short URL in http://localhost:8080/id format.
//
// Encodes given URL using base64 encoding scheme and puts it to the URL's map.
//
// Returns 201 status code if processing was successfull, otherwise returns 400.
func PostHandlerURL(writer http.ResponseWriter, req *http.Request) {
	logger.Log.Debug("POST handler URL processing")

	inputURL, err := io.ReadAll(req.Body)
	defer req.Body.Close()

	if err != nil {
		logger.Log.Error(CannotProcessURL, zap.Error(err))
		http.Error(writer, CannotProcessURL, http.StatusBadRequest)
		return
	}

	if ok := entity.IsValidURL(string(inputURL)); !ok {
		logger.Log.Error(WrongURLFormat, zap.Error(err))
		http.Error(writer, WrongURLFormat, http.StatusBadRequest)
		return
	}

	baseURIPrefix := req.Context().Value(baseURIPrefixCtx).(string)
	if baseURIPrefix == "" {
		logger.Log.Error("invalid base URI prefix", zap.String("base URI prefix", baseURIPrefix))
		http.Error(writer, InternalServerError, http.StatusInternalServerError)
		return
	}

	outputURL := postURLProcessing(string(inputURL), baseURIPrefix)
	if outputURL == "" {
		logger.Log.Error("could not create a short URL")
		http.Error(writer, InternalServerError, http.StatusInternalServerError)
		return
	}

	logger.Log.Info("url has been created succeessfully", zap.String("output url", outputURL))

	writer.Header().Set("Content-Type", "text/plain; charset=utf-8")
	writer.WriteHeader(http.StatusCreated)
	io.WriteString(writer, outputURL)
}

func PostHandlerJSON(writer http.ResponseWriter, req *http.Request) {
	logger.Log.Debug("POST handler JSON processing")

	inputRequest := &models.Request{}
	err := json.NewDecoder(req.Body).Decode(&inputRequest)
	if err != nil {
		logger.Log.Error(CannotProcessJSON, zap.Error(err))
		http.Error(writer, CannotProcessJSON, http.StatusBadRequest)
		return
	}

	if ok := entity.IsValidURL(inputRequest.URL); !ok {
		logger.Log.Error(WrongURLFormat, zap.Error(err))
		http.Error(writer, WrongURLFormat, http.StatusBadRequest)
		return
	}

	baseURIPrefix := req.Context().Value(baseURIPrefixCtx).(string)
	if baseURIPrefix == "" {
		logger.Log.Error("invalid base URI prefix", zap.String("base URI prefix", baseURIPrefix))
		http.Error(writer, InternalServerError, http.StatusInternalServerError)
		return
	}

	outputURL := postURLProcessing(inputRequest.URL, baseURIPrefix)
	if outputURL == "" {
		logger.Log.Error("could not create a short URL")
		http.Error(writer, InternalServerError, http.StatusInternalServerError)
		return
	}

	response := &models.Response{
		URL: outputURL,
	}
	logger.Log.Info("url has been created succeessfully", zap.String("output url", response.URL))

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusCreated)
	if err = json.NewEncoder(writer).Encode(response); err != nil {
		logger.Log.Error("invalid response", zap.Any("response", response))
		http.Error(writer, "internal server error", http.StatusInternalServerError)
		return
	}

	logger.Log.Debug("sending HTTP 200 response")
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
