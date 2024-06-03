package handlers

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

// URLHandler Processes POST "/id" endpoint. Save original and short URLs to storage
//
// Returns 201(StatusCreated) if processing was successfully
// Returns 500(StatusInternalServerError) if base URI prefix is invalid
// Returns 500(StatusInternalServerError) if user ID is invalid
// Returns 500(StatusInternalServerError) when database error
// Returns 400(StatusBadRequest) if original URL is invalid
// Returns 409(StatusConflict) if original URL exists in storage for this user
func URLHandler(saver URLSaver, baseURIPrefix string) http.HandlerFunc {
	return func(writer http.ResponseWriter, req *http.Request) {
		zap.L().Debug("POST handler URL processing")

		if baseURIPrefix == "" {
			zap.L().Error("invalid base URI prefix", zap.String("base URI prefix", baseURIPrefix))
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}

		userIDCtx, ok := req.Context().Value(entity.UserIDCtxKey{}).(entity.UserIDCtx)
		if !ok {
			zap.L().Error("user id couldn't obtain from context")
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}

		if len(userIDCtx.UserID.String()) == 0 {
			zap.L().Error("empty user id from context while posting user url")
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}

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

		ctx, cancel := context.WithTimeout(req.Context(), timeout)
		defer cancel()

		outputURL, err := postURLProcessing(saver, ctx, userIDCtx.UserID, string(inputURL), baseURIPrefix)
		if err != nil {
			zap.L().Error("could not create a short URL", zap.String("error", err.Error()))
			if errors.Is(err, storage_err.ErrURLAlreadyExists) {
				successRawResponse(writer, outputURL, http.StatusConflict)
				return
			}

			writer.WriteHeader(http.StatusInternalServerError)
			return
		}

		zap.L().Info("url has been created successfully", zap.String("output url", outputURL))

		successRawResponse(writer, outputURL, http.StatusCreated)
	}
}

func successRawResponse(writer http.ResponseWriter, url string, status int) {
	writer.Header().Set("Content-Type", "text/plain; charset=utf-8")
	writer.WriteHeader(status)
	io.WriteString(writer, url)
}
