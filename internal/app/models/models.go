package models

// Request Contains information about original URL in JSON representation
type Request struct {
	URL string `json:"url"`
}

// Response Contains information about short URL in JSON representation
type Response struct {
	URL string `json:"result"`
}
