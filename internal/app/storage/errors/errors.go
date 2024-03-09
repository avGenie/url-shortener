package errors

import "errors"

var (
	// ErrEmptyDBStorageCred = errors.New("error in db storage connection credentials are empty")
	ErrShortURLNotFound = errors.New("short url is not found in storage")
)