package grpc

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/avGenie/url-shortener/internal/app/entity"
	"github.com/avGenie/url-shortener/internal/app/grpc/converter"
	grpc_context "github.com/avGenie/url-shortener/internal/app/grpc/usecase/context"
	get_handlers "github.com/avGenie/url-shortener/internal/app/handlers/get"
	post_handlers "github.com/avGenie/url-shortener/internal/app/handlers/post"
	"github.com/avGenie/url-shortener/internal/app/models"
	storage_err "github.com/avGenie/url-shortener/internal/app/storage/api/errors"
	pb "github.com/avGenie/url-shortener/proto"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

const (
	requestTimeout = 3 * time.Second

	ErrInternalMsg           = "internal server error"
	ErrEmptyBaseURIPrefixMsg = "base uri prefix is empty"
	ErrWrongURLFormatMsg     = "wrong URL format"
)

// GetShortURL Returns short URL by original and user id
func (s *ShortenerServer) GetShortURL(ctx context.Context, original *pb.OriginalURL) (*pb.ShortURL, error) {
	userID := grpc_context.GetUserIDFromContext(ctx)

	if s.config.BaseURIPrefix == "" {
		zap.L().Error(ErrEmptyBaseURIPrefixMsg)

		return nil, status.Errorf(codes.Internal, ErrInternalMsg)
	}

	if ok := entity.IsValidURL(original.Url); !ok {
		zap.L().Error(ErrWrongURLFormatMsg)

		return nil, status.Errorf(codes.InvalidArgument, "couldn't parse %s url", original)
	}

	ctx, cancel := context.WithTimeout(ctx, requestTimeout)
	defer cancel()

	outputURL, err := post_handlers.PostURLProcessing(
		s.storage,
		ctx,
		userID,
		original.Url,
		s.config.BaseURIPrefix,
	)
	if err != nil {
		zap.L().Error("could not create a short URL", zap.String("error", err.Error()))
		if errors.Is(err, storage_err.ErrURLAlreadyExists) {
			return nil, status.Errorf(codes.AlreadyExists, "url already exists in storage for this user")
		}

		return nil, status.Errorf(codes.Internal, ErrInternalMsg)
	}

	return &pb.ShortURL{Url: outputURL}, nil
}

// GetOriginalURL Returns original URL by short and user id
func (s *ShortenerServer) GetOriginalURL(ctx context.Context, original *pb.ShortURL) (*pb.OriginalURL, error) {
	userID := grpc_context.GetUserIDFromContext(ctx)

	shortURL, err := entity.ParseURL(original.Url)
	if err != nil {
		zap.L().Error("couldn't parse original URL", zap.Error(err), zap.String("user_id", userID.String()))

		return nil, status.Errorf(codes.InvalidArgument, "couldn't parse %s url", original)
	}

	ctx, cancel := context.WithTimeout(ctx, requestTimeout)
	defer cancel()

	url, err := s.storage.GetURL(ctx, userID, *shortURL)
	if err != nil {
		if errors.Is(err, storage_err.ErrAllURLsDeleted) {
			errMsg := "original url has been deleted for this user"
			zap.L().Error(errMsg, zap.Error(err), zap.String("user_id", userID.String()))

			return nil, status.Errorf(codes.NotFound, errMsg)
		}

		zap.L().Error(
			"error while getting url",
			zap.String("error", err.Error()),
			zap.String("short_url", shortURL.String()),
		)

		return nil, status.Errorf(codes.Internal, ErrInternalMsg)
	}

	return &pb.OriginalURL{Url: url.String()}, nil
}

// GetAllUserURL Returns all user URLs
func (s *ShortenerServer) GetAllUserURL(ctx context.Context, _ *emptypb.Empty) (*pb.AllUrlsResponse, error) {
	userID := grpc_context.GetUserIDFromContext(ctx)

	if s.config.BaseURIPrefix == "" {
		zap.L().Error(ErrEmptyBaseURIPrefixMsg)

		return nil, status.Errorf(codes.Internal, ErrInternalMsg)
	}

	ctx, cancel := context.WithTimeout(ctx, requestTimeout)
	defer cancel()

	data, err := get_handlers.ProcessAllUserURL(s.storage, ctx, userID, s.config.BaseURIPrefix)
	if err != nil {
		if errors.Is(err, get_handlers.ErrAllURLNotFound) {
			errMsg := "all urls not found for given user"
			zap.L().Error(errMsg, zap.String("user_id", userID.String()))

			return nil, status.Errorf(codes.NotFound, errMsg)
		}

		zap.L().Error("error while processing all user urls", zap.Error(err))

		return nil, status.Errorf(codes.Internal, ErrInternalMsg)
	}

	var batch models.AllUrlsBatch
	err = json.Unmarshal(data, &batch)
	if err != nil {
		zap.L().Error("error while unmarshalling all user data", zap.Error(err))

		return nil, status.Errorf(codes.Internal, ErrInternalMsg)
	}

	urlsResponse := converter.AllURLsBatchToAllURLsResponse(batch)

	return urlsResponse, nil
}

// GetBatchShortURL Returns short URL by original batch of URLs and user id
func (s *ShortenerServer) GetBatchShortURL(ctx context.Context, originalBatch *pb.BatchRequest) (*pb.BatchResponse, error) {
	userID := grpc_context.GetUserIDFromContext(ctx)

	if s.config.BaseURIPrefix == "" {
		zap.L().Error(ErrEmptyBaseURIPrefixMsg)

		return nil, status.Errorf(codes.Internal, ErrInternalMsg)
	}

	reqBatch := converter.BatchRequestToReqBatch(originalBatch)

	ctx, cancel := context.WithTimeout(ctx, requestTimeout)
	defer cancel()

	resBatch, err := post_handlers.BatchURLProcessing(s.storage, ctx, userID, reqBatch, s.config.BaseURIPrefix)
	if err != nil {
		zap.L().Error("error while batch url processing", zap.Error(err))

		return nil, status.Errorf(codes.Internal, ErrInternalMsg)
	}

	outBatch := converter.ResBatchToBatchResponse(resBatch)

	return outBatch, nil
}
