package post

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/avGenie/url-shortener/internal/app/converter"
	"github.com/avGenie/url-shortener/internal/app/entity"
	post_err "github.com/avGenie/url-shortener/internal/app/handlers/errors"
	"github.com/avGenie/url-shortener/internal/app/models"
	storage_err "github.com/avGenie/url-shortener/internal/app/storage/api/errors"
	"go.uber.org/zap"
)

// Processes POST request by JSON within http://localhost:8080/api/shorten URL format.
//
// Encodes given URL using base64 encoding scheme and puts it to the URL's map.
//
// Returns 201 status code if processing was successfull, otherwise returns 400.
func JSONHandler(saver URLSaver, baseURIPrefix string) http.HandlerFunc {
	return func(writer http.ResponseWriter, req *http.Request) {
		zap.L().Debug("POST handler JSON processing")

		inputRequest := &models.Request{}
		err := json.NewDecoder(req.Body).Decode(&inputRequest)
		defer req.Body.Close()
		if err != nil {
			zap.L().Error(post_err.CannotProcessJSON, zap.Error(err))
			http.Error(writer, post_err.WrongJSONFormat, http.StatusBadRequest)
			return
		}

		if ok := entity.IsValidURL(inputRequest.URL); !ok {
			zap.L().Error(post_err.WrongJSONFormat, zap.Error(err))
			http.Error(writer, post_err.WrongJSONFormat, http.StatusBadRequest)
			return
		}

		if baseURIPrefix == "" {
			zap.L().Error("invalid base URI prefix", zap.String("base URI prefix", baseURIPrefix))
			http.Error(writer, post_err.InternalServerError, http.StatusInternalServerError)
			return
		}

		successResp := func(response models.Response, status int) {
			zap.L().Info("url has been created succeessfully", zap.String("output url", response.URL))

			writer.Header().Set("Content-Type", "application/json")
			writer.WriteHeader(status)
			if err = json.NewEncoder(writer).Encode(response); err != nil {
				zap.L().Error("invalid response", zap.Any("response", response))
				http.Error(writer, "internal server error", http.StatusInternalServerError)
				return
			}
		}

		ctx, cancel := context.WithTimeout(req.Context(), timeout)
		defer cancel()

		outputURL, err := postURLProcessing(saver, ctx, inputRequest.URL, baseURIPrefix)

		response := models.Response{
			URL: outputURL,
		}

		if err != nil {
			zap.L().Error("could not create a short URL", zap.String("error", err.Error()))
			if errors.Is(err, storage_err.ErrURLAlreadyExists) {
				successResp(response, http.StatusConflict)
				return
			}

			http.Error(writer, post_err.InternalServerError, http.StatusInternalServerError)
			return
		}

		successResp(response, http.StatusCreated)

		zap.L().Debug("sending HTTP 200 response")
	}
}

func JSONBatchHandler(saver URLBatchSaver, baseURIPrefix string) http.HandlerFunc {
	return func(writer http.ResponseWriter, req *http.Request) {
		zap.L().Debug("POST JSON batch handler processing")

		var batch models.ReqBatch
		err := json.NewDecoder(req.Body).Decode(&batch)
		if err != nil {
			zap.L().Error(post_err.CannotProcessJSON, zap.Error(err))
			http.Error(writer, post_err.WrongJSONFormat, http.StatusBadRequest)
			return
		}
		defer req.Body.Close()

		urls, err := converter.ConvertBatchReqToURL(batch)
		if err != nil {
			zap.L().Error(post_err.CannotProcessURL, zap.Error(err))
			http.Error(writer, post_err.WrongJSONFormat, http.StatusBadRequest)
			return
		}

		sBatch, err := createStorageBatch(urls)
		if err != nil {
			zap.L().Error("error while creating storage batch", zap.Error(err))
			http.Error(writer, post_err.InternalServerError, http.StatusInternalServerError)
			return
		}

		ctx, cancel := context.WithTimeout(req.Context(), timeout)
		defer cancel()

		savedBatch, err := saver.SaveBatchURL(ctx, sBatch)
		if err != nil {
			zap.L().Error("error while saving url to storage", zap.Error(err))
			http.Error(writer, post_err.InternalServerError, http.StatusInternalServerError)
			return
		}

		outBatch := converter.ConvertStorageBatchToOutBatch(savedBatch, baseURIPrefix)
		out, err := json.Marshal(outBatch)
		if err != nil {
			zap.L().Error("error while converting storage url to output", zap.Error(err))
			http.Error(writer, post_err.InternalServerError, http.StatusInternalServerError)
			return
		}

		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusCreated)
		writer.Write(out)
	}
}
