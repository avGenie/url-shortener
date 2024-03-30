package api

import "errors"

var (
	ErrShortURLNotFound   = errors.New("short url is not found in storage for this user")
	ErrURLAlreadyExists   = errors.New("short url already exists in storage for this user")
	ErrFileStorageNotOpen = errors.New("file storage is not open")
	ErrUserIDNotFound     = errors.New("user id not found in storage")
)
