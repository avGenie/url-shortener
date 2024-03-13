package post

import (
	"context"
	"errors"
	"io"
	"net/http"

	"github.com/avGenie/url-shortener/internal/app/entity"
	post_err "github.com/avGenie/url-shortener/internal/app/handlers/errors"
	storage_err "github.com/avGenie/url-shortener/internal/app/storage/api/errors"
	"go.uber.org/zap"
)

// Processes POST request by URL within http://localhost:8080/id URL format.
//
// Encodes given URL using base64 encoding scheme and puts it to the URL's map.
//
// Returns 201 status code if processing was successfull, otherwise returns 400.
func URLHandler(saver URLSaver, baseURIPrefix string) http.HandlerFunc {
	return func(writer http.ResponseWriter, req *http.Request) {
		zap.L().Debug("POST handler URL processing")

		inputURL, err := io.ReadAll(req.Body)
		defer req.Body.Close()

		if err != nil {
			zap.L().Error(post_err.CannotProcessURL, zap.Error(err))
			http.Error(writer, post_err.WrongURLFormat, http.StatusBadRequest)
			return
		}

		if ok := entity.IsValidURL(string(inputURL)); !ok {
			zap.L().Error(post_err.WrongURLFormat, zap.Error(err))
			http.Error(writer, post_err.WrongURLFormat, http.StatusBadRequest)
			return
		}

		if baseURIPrefix == "" {
			zap.L().Error("invalid base URI prefix", zap.String("base URI prefix", baseURIPrefix))
			http.Error(writer, post_err.InternalServerError, http.StatusInternalServerError)
			return
		}

		successResp := func(url string, status int) {
			writer.Header().Set("Content-Type", "text/plain; charset=utf-8")
			writer.WriteHeader(status)
			io.WriteString(writer, url)
		}

		ctx, cancel := context.WithTimeout(req.Context(), timeout)
		defer cancel()

		outputURL, err := postURLProcessing(saver, ctx, string(inputURL), baseURIPrefix)
		if err != nil {
			zap.L().Error("could not create a short URL", zap.String("error", err.Error()))
			if errors.Is(err, storage_err.ErrURLAlreadyExists) {
				successResp(outputURL, http.StatusConflict)
				return
			}

			http.Error(writer, post_err.InternalServerError, http.StatusInternalServerError)
			return
		}

		zap.L().Info("url has been created succeessfully", zap.String("output url", outputURL))

		successResp(outputURL, http.StatusCreated)
	}
}
