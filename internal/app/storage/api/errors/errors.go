package api

import "errors"

// Errors which return storage
var (
	ErrShortURLNotFound   = errors.New("short url is not found in storage for this user")
	ErrURLAlreadyExists   = errors.New("short url already exists in storage for this user")
	ErrFileStorageNotOpen = errors.New("file storage is not open")
	ErrAllURLsDeleted     = errors.New("all urls have been deleted for this user")
)
