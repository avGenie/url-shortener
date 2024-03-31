package entity

type URLRecord struct {
	ID          uint   `json:"uuid"`
	UserID      string `json:"user_id"`
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}
