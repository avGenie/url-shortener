package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"time"

	"github.com/avGenie/url-shortener/internal/app/entity"
	handler_err "github.com/avGenie/url-shortener/internal/app/handlers/errors"
	"github.com/avGenie/url-shortener/internal/app/models"
	"go.uber.org/zap"
)

const (
	tickerTime  = 5 * time.Second
	contextTime = 3 * time.Second
)

type AllURLDeleter interface {
	DeleteBatchURL(ctx context.Context, urls entity.DeletedURLBatch) error
}

type DeleteHandler struct {
	deleter AllURLDeleter

	msgChan chan []entity.DeletedURL
}

func NewDeleteHandler(deleter AllURLDeleter) *DeleteHandler {
	instance := &DeleteHandler{
		deleter: deleter,
		msgChan: make(chan []entity.DeletedURL),
	}

	go instance.flushDeletedURLs()

	return instance
}

func (h *DeleteHandler) DeleteUserURLHandler() http.HandlerFunc {
	return func(writer http.ResponseWriter, req *http.Request) {
		userIDCtx, ok := req.Context().Value(entity.UserIDCtxKey{}).(entity.UserIDCtx)
		if !ok {
			zap.L().Error("user id couldn't obtain from context while all user urls deleting")
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}

		if len(userIDCtx.UserID.String()) == 0 {
			zap.L().Error("empty user id from context while posting all user urls deleting")
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
	for urls := range h.msgChan {
		ctx, cancel := context.WithTimeout(context.Background(), contextTime)

		err := h.deleter.DeleteBatchURL(ctx, urls)
		if err != nil {
			switch {
			case errors.Is(err, context.Canceled):
				zap.L().Error("context canceled while flushing deleted urls", zap.String("error", err.Error()))
			case errors.Is(err, context.DeadlineExceeded):
				zap.L().Error("context deadline exceeded while flushing deleted urls", zap.String("error", err.Error()))
			default:
				zap.L().Error("error while flushing deleted urls", zap.Error(err))
			}
		}

		cancel()
	}
}

func (h *DeleteHandler) processDeletedURLs(userID entity.UserID, reader io.ReadCloser) error {
	var batch models.ReqDeletedURLBatch
	err := json.NewDecoder(reader).Decode(&batch)
	if err != nil {
		zap.L().Error("cannot process input user urls for deleting", zap.Error(err))
		return handler_err.ErrWrongDeletedURLFormat
	}
	defer reader.Close()

	resURLBatch := make([]entity.DeletedURL, 0, len(batch))
	for _, url := range batch {
		resURLBatch = append(resURLBatch, entity.DeletedURL{
			UserID:   userID.String(),
			ShortURL: url,
		})
	}

	h.msgChan <- resURLBatch

	return nil
}
