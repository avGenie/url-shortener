package errors

import "errors"

// Handler processing errors
//
// WrongURLFormat - returned as HTTP output if an invalid URL format was received
// WrongJSONFormat - returned as HTTP output if an invalid URL JSON format was received
// ShortURLNotInDB - returned as HTTP output if given short URL is not found in DB
// CannotProcessURL - returned as log message if URL couldn't be processed
// CannotProcessJSON - returned as log message if URL couldn't be processed in JSON format
// InternalServerError - returned as HTTP output due to internal server error
// ErrWrongDeletedURLFormat - returned if the url could not be deleted
var (
	WrongURLFormat    = "wrong URL format"
	WrongJSONFormat   = "wrong JSON format"
	ShortURLNotInDB   = "given short URL did not find in database"
	CannotProcessURL  = "cannot process URL"
	CannotProcessJSON = "cannot process JSON"

	InternalServerError = "internal server error"

	ErrWrongDeletedURLFormat = errors.New("wrong deleted urls format")
)
