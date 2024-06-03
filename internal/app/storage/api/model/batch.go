package model

// Batch Slice of BatchObject objects
type Batch []BatchObject

// BatchObject Contains information about batch URL
type BatchObject struct {
	ID       string
	InputURL string
	ShortURL string
}
