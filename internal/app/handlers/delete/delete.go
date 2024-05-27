// Package handlers provides application endpoints
package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"sync"
	"time"

	"github.com/avGenie/url-shortener/internal/app/entity"
	handler_err "github.com/avGenie/url-shortener/internal/app/handlers/errors"
	"github.com/avGenie/url-shortener/internal/app/models"
	"go.uber.org/zap"
)

const (
	flushBufLen = 100

	tickerTime  = 5 * time.Second
	contextTime = 3 * time.Second
	stopTimeout = 5 * time.Second
)

// AllURLDeleter Storage interface for delete handler
type AllURLDeleter interface {
	DeleteBatchURL(ctx context.Context, urls entity.DeletedURLBatch) error
}

// DeleteHandler Endpoints for delete operations
type DeleteHandler struct {
	deleter AllURLDeleter

	wg      *sync.WaitGroup
	done    chan struct{}
	msgChan chan entity.DeletedURLBatch
}

// NewDeleteHandler Creates delete handler using obtained storage
func NewDeleteHandler(deleter AllURLDeleter) *DeleteHandler {
	instance := &DeleteHandler{
		deleter: deleter,
		wg:      &sync.WaitGroup{},
		done:    make(chan struct{}),
		msgChan: make(chan entity.DeletedURLBatch),
	}

	instance.wg.Add(1)
	go func() {
		defer instance.wg.Done()
		instance.flushDeletedURLs()
	}()

	return instance
}

// DeleteUserURLHandler Process endpoint for user deletion
//
// Returns 202(StatusAccepted) if deletion was successful
// Returns 500(StatusInternalServerError) if user id is incorrect
// Returns 500(StatusInternalServerError) if user id could not be parsed
// Returns 400(StatusBadRequest) if user id could not be processed for deletion
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

// Stop Stops user deletion process
func (h *DeleteHandler) Stop() {
	sync.OnceFunc(func() {
		close(h.done)
	})()

	ready := make(chan bool)
	go func() {
		defer close(ready)
		h.wg.Wait()
	}()

	// устанавливаем таймаут на ожидание сброса в БД последней порции
	select {
	case <-time.After(stopTimeout):
		zap.L().Error("timeout stopped while sending data to the storage while shutting down")
		return
	case <-ready:
		zap.L().Info("succsessful sending data to the storage while shutting down")
		return
	}
}

func (h *DeleteHandler) flushDeletedURLs() {
	ticker := time.NewTicker(tickerTime)
	storageBatch := make([]entity.DeletedURL, 0, flushBufLen)

	flush := func() {
		ctx, cancel := context.WithTimeout(context.Background(), contextTime)
		defer cancel()

		zap.L().Debug("flushing deleted urls", zap.Int("urls_count", len(storageBatch)))

		err := h.deleter.DeleteBatchURL(ctx, storageBatch)
		if err != nil {
			switch {
			case errors.Is(err, context.Canceled):
				zap.L().Error("context canceled while flushing deleted urls", zap.String("error", err.Error()))
			case errors.Is(err, context.DeadlineExceeded):
				zap.L().Error("context deadline exceeded while flushing deleted urls", zap.String("error", err.Error()))
			default:
				zap.L().Error("error while flushing deleted urls", zap.Error(err))
			}

			return
		}

		storageBatch = storageBatch[:0:flushBufLen]
	}

	for {
		select {
		case <-h.done:
			zap.L().Info("shutting down server; last flushing")
			flush()
			return
		case urls, ok := <-h.msgChan:
			if !ok {
				return
			}

			if len(storageBatch)+len(urls) > flushBufLen {
				flush()
			}

			storageBatch = append(storageBatch, urls...)
		case <-ticker.C:
			if len(storageBatch) == 0 {
				continue
			}

			zap.L().Debug("flushing deleted urls by ticker")
			flush()
		}
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
