// Package entity contains structs which are used in application
package entity

// DeletedURLBatch Array of DeletedURL structs
type DeletedURLBatch []DeletedURL

// DeletedURL is being used to generate data to be sent to the database for deletion
type DeletedURL struct {
	UserID   string
	ShortURL string
}
