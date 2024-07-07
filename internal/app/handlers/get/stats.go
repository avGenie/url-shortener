package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/avGenie/url-shortener/internal/app/models"
	cidr "github.com/avGenie/url-shortener/internal/app/usecase/CIDR"
	"go.uber.org/zap"
)

const realIPKey = "X-Real-IP"

// StatisticGetter Getter for service statistic request
type StatisticGetter interface {
	GetStatistic(ctx context.Context) (models.CountStatistic, error)
}

// StatsHandler Processes service statistic request
//
// Returns 200(StatusOk) if processing was successful
// Returns 500(StatusInternalServerError) when parsing or DB request errors
// Returns 403(StatusForbidden) when request forbidden for given IP
func StatsHandler(statGetter StatisticGetter, cidr *cidr.CIDR) http.HandlerFunc {
	return func(writer http.ResponseWriter, req *http.Request) {
		zap.L().Debug("stats handler URL processing")

		err := processCIDR(req, cidr)
		if err != nil {
			zap.L().Info("forbidden to get statistic", zap.Error(err))

			writer.WriteHeader(http.StatusForbidden)

			return
		}

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

func processCIDR(req *http.Request, cidr *cidr.CIDR) error {
	if cidr == nil {
		return fmt.Errorf("subnet unknown")
	}

	userIP := req.Header.Get(realIPKey)
	isSubnet := cidr.Contains(userIP)
	if !isSubnet {
		return fmt.Errorf("user ip is not in subnet")
	}

	return nil
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
