package model

import "github.com/avGenie/url-shortener/internal/app/entity"

type Batch []BatchObject

type BatchObject struct {
	ID       string
	InputURL string
	ShortURL string
}

type BatchResponse struct {
	entity.Response

	Batch Batch
}

func OKBatchResponse(batch []BatchObject) BatchResponse {
	return BatchResponse{
		Response: entity.OKResponse(),
		Batch:    batch,
	}
}

func ErrorBatchResponse(err error) BatchResponse {
	return BatchResponse{
		Response: entity.ErrorResponse(err),
		Batch:    nil,
	}
}
