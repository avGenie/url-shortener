package api

import "errors"

// Errors which return storage
//
// ErrShortURLNotFound - returned if URL is not found in storage for user
// ErrURLAlreadyExists - returned if short URL already exists in storage for user
// ErrFileStorageNotOpen - returned if file storage is not opened
// ErrAllURLsDeleted - returned if all URLs deleted for user
var (
	ErrShortURLNotFound   = errors.New("short url is not found in storage for this user")
	ErrURLAlreadyExists   = errors.New("short url already exists in storage for this user")
	ErrFileStorageNotOpen = errors.New("file storage is not open")
	ErrAllURLsDeleted     = errors.New("all urls have been deleted for this user")
)
