package models

import "github.com/avGenie/url-shortener/internal/app/entity"

type ReqBatch []BatchObjectRequest
type ReqURLBatch []ReqURLBatchObject
type ResBatch []BatchObjectResponse

type ReqURLBatchObject struct {
	Obj BatchObjectRequest
	URL entity.URL
}

type BatchObjectRequest struct {
	ID  string `json:"correlation_id"`
	URL string `json:"original_url"`
}

type BatchObjectResponse struct {
	ID  string `json:"correlation_id"`
	URL string `json:"short_url"`
}
