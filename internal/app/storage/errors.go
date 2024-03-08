package storage

import "errors"

var (
	ErrEmptyDBStorageCred = errors.New("error in db storage connection credentials are empty")
)