package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/avGenie/url-shortener/internal/app/models"
	"go.uber.org/zap"
)

type StatisticGetter interface {
	GetStatistic(ctx context.Context) (models.CountStatistic, error)
}

func StatsHandler(statGetter StatisticGetter) http.HandlerFunc {
	return func(writer http.ResponseWriter, req *http.Request) {
		zap.L().Debug("stats handler URL processing")

		out, err := processServiceStatistic(req.Context(), statGetter)
		if err != nil {
			writer.WriteHeader(http.StatusInternalServerError)

			return
		}

		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusOK)
		writer.Write(out)
	}
}

func processServiceStatistic(ctx context.Context, statGetter StatisticGetter) ([]byte, error) {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	stat, err := statGetter.GetStatistic(ctx)
	if err != nil {
		zap.L().Error(
			"error while getting statistic",
			zap.String("error", err.Error()),
		)

		return nil, err
	}

	out, err := json.Marshal(stat)
	if err != nil {
		errMsg := "error while converting service statistic to output"
		zap.L().Error(errMsg, zap.Error(err))

		return nil, fmt.Errorf("%s: %w", errMsg, err)
	}

	return out, nil
}
