package handlers

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/avGenie/url-shortener/internal/app/entity"
	"github.com/avGenie/url-shortener/internal/app/handlers/errors"
	"github.com/avGenie/url-shortener/internal/app/models"
	storage "github.com/avGenie/url-shortener/internal/app/storage/api/model"
	"go.uber.org/zap"
)

const (
	tickerTime = 5 * time.Second
	contextTime = 3 * time.Second
)

type DeleteHandler struct {
	storage storage.Storage

	msgChan chan entity.DeletedURL
}

func NewDeleteHandler(storage storage.Storage) *DeleteHandler {
	instance := &DeleteHandler{
		storage: storage,
		msgChan: make(chan entity.DeletedURL, 100),
	}

	go instance.flushDeletedURLs()

	return instance
}

func (h *DeleteHandler) DeleteUserURLHandler() http.HandlerFunc {
	return func(writer http.ResponseWriter, req *http.Request) {
		userIDCtx, ok := req.Context().Value(entity.UserIDCtxKey{}).(entity.UserIDCtx)
		if !ok {
			zap.L().Error("user id couldn't obtain from context while all user urls processing")
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}

		err := h.processDeletedURLs(userIDCtx.UserID, req.Body)
		if err != nil {
			writer.WriteHeader(http.StatusBadRequest)
			return
		}

		writer.WriteHeader(http.StatusAccepted)
	}
}

func (h *DeleteHandler) flushDeletedURLs() {
	ticker := time.NewTicker(tickerTime)

	var storageBatch entity.DeletedURLBatch
	for {
		select {
		case url, ok := <- h.msgChan:
			if !ok {
				return
			}

			storageBatch = append(storageBatch, url)
		case <- ticker.C:
			if len(storageBatch) == 0 {
				continue
			}

			ctx, cancel := context.WithTimeout(context.Background(), contextTime)
			defer cancel()

			err := h.storage.DeleteBatchURL(ctx, storageBatch)
			if err != nil {
				zap.L().Error("cannot delete urls", zap.Error(err))
				continue
			}

			storageBatch = nil
		}
	}
}

func (h *DeleteHandler) processDeletedURLs(userID entity.UserID, reader io.ReadCloser) error {
	var batch models.ReqDeletedURLBatch
	err := json.NewDecoder(reader).Decode(&batch)
	if err != nil {
		zap.L().Error("cannot process input user urls for deleting", zap.Error(err))
		return errors.ErrWrongDeletedURLFormat
	}
	defer reader.Close()

	for _, url := range batch {
		h.msgChan <- entity.DeletedURL{
			UserID:   userID.String(),
			ShortURL: url,
		}
	}

	return nil
}
