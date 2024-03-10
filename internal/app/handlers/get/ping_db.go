package get

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/avGenie/url-shortener/internal/app/entity"
	"go.uber.org/zap"
)

const (
	pingTimeout = 1*time.Second
)

type StoragePinger interface {
	PingServer(ctx context.Context) entity.Response
}

func GetPingDB(pinger StoragePinger) http.HandlerFunc {
	return func(writer http.ResponseWriter, req *http.Request) {
		ctx, cancel := context.WithTimeout(req.Context(), pingTimeout)
		defer cancel()

		response := pinger.PingServer(ctx)
		if response.Status == entity.StatusError {
			switch {
			case errors.Is(response.Error, context.Canceled):
				zap.L().Error("context canceled", zap.String("error", response.Error.Error()))
			case errors.Is(response.Error, context.DeadlineExceeded):
				zap.L().Error("context deadline exceeded", zap.String("error", response.Error.Error()))
			default:
				zap.L().Error(response.Error.Error())
			}
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}

		writer.WriteHeader(http.StatusOK)
	}
}