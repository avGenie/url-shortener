package handlers

import (
	"context"
	"errors"
	"net/http"

	"go.uber.org/zap"
)

// StoragePinger Interface to ping storage
type StoragePinger interface {
	PingServer(ctx context.Context) error
}

// PingDBHandler Processes GET "/ping" endpoint. Sends ping to storage
//
// Returns 200(StatusOK) if processing was successful
// Returns 500(StatusInternalServerError) if the database could not be accessed
func PingDBHandler(pinger StoragePinger) http.HandlerFunc {
	return func(writer http.ResponseWriter, req *http.Request) {
		ctx, cancel := context.WithTimeout(req.Context(), timeout)
		defer cancel()

		err := pinger.PingServer(ctx)
		if err != nil {
			switch {
			case errors.Is(err, context.Canceled):
				zap.L().Error("context canceled", zap.String("error", err.Error()))
			case errors.Is(err, context.DeadlineExceeded):
				zap.L().Error("context deadline exceeded", zap.String("error", err.Error()))
			default:
				zap.L().Error(err.Error())
			}
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}

		zap.L().Info("storage works after ping")

		writer.WriteHeader(http.StatusOK)
	}
}
