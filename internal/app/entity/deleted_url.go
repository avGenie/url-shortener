package entity

type DeletedURLBatch []DeletedURL

type DeletedURL struct {
	UserID   string
	ShortURL string
}
