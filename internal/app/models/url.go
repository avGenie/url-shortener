package models

type AllUrlsBatch []AllUrlsResponse

type AllUrlsResponse struct {
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}
