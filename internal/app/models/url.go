package models

// AllUrlsBatch Slice of AllUrlsResponse structs
type AllUrlsBatch []AllUrlsResponse

// AllUrlsResponse Contains information about original and short URL in JSON representation
type AllUrlsResponse struct {
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}
