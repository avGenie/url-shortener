// Package models contains structs which are used for communication with external services
package models

import "github.com/avGenie/url-shortener/internal/app/entity"

// ReqBatch Slice of structs for passing batch URL data to storage
type ReqBatch []BatchObjectRequest

// ReqURLBatch Input slice of structs for batch POST request
type ReqURLBatch []ReqURLBatchObject

// ResBatch Output slice of structs for batch POST request
type ResBatch []BatchObjectResponse

// ReqURLBatchObject Struct for passing batch URL data to storage
type ReqURLBatchObject struct {
	Obj BatchObjectRequest
	URL entity.URL
}

// BatchObjectRequest Input struct for batch POST request
type BatchObjectRequest struct {
	ID  string `json:"correlation_id"`
	URL string `json:"original_url"`
}

// BatchObjectResponse Output struct for batch POST request
type BatchObjectResponse struct {
	ID  string `json:"correlation_id"`
	URL string `json:"short_url"`
}
