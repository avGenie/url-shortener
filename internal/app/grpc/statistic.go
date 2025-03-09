package grpc

import (
	"context"

	"github.com/avGenie/url-shortener/internal/app/grpc/converter"
	get_handlers "github.com/avGenie/url-shortener/internal/app/handlers/get"
	pb "github.com/avGenie/url-shortener/proto"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

// GetStatistic Returns server statistic
func (s *ShortenerServer) GetStatistic(ctx context.Context, _ *emptypb.Empty) (*pb.StatisticResposne, error) {
	ctx, cancel := context.WithTimeout(ctx, requestTimeout)
	defer cancel()

	stat, err := get_handlers.ProcessServiceStatistic(ctx, s.storage)
	if err != nil {
		zap.L().Error("could not process statistic", zap.String("error", err.Error()))

		return nil, status.Errorf(codes.Internal, ErrInternalMsg)
	}

	output := converter.CountStatisticToStatisticResposne(stat)

	return output, nil
}
