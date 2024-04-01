package errors

import "errors"

var (
	WrongURLFormat    = "wrong URL format"
	WrongJSONFormat   = "wrong JSON format"
	ShortURLNotInDB   = "given short URL did not find in database"
	CannotProcessURL  = "cannot process URL"
	CannotProcessJSON = "cannot process JSON"

	InternalServerError = "internal server error"

	ErrWrongDeletedURLFormat = errors.New("wrong deleted urls format")
)
