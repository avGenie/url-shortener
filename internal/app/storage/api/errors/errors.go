package api

import "errors"

var (
	ErrShortURLNotFound   = errors.New("short url is not found in storage")
	ErrURLAlreadyExists   = errors.New("short url already exists in storage")
	ErrFileStorageNotOpen = errors.New("file storage is not open")
)
