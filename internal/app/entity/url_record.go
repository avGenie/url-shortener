package entity

// URLRecord is being used to form a string for the file database
type URLRecord struct {
	ID          uint   `json:"uuid"`
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}
