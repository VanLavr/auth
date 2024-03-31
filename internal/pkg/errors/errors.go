package errors

import "errors"

var (
	ErrTokenWasNotProvided  = errors.New("auth token was not provided")
	ErrInvalidSigningMethod = errors.New("provided token was signed via invalid method")
	ErrInternal             = errors.New("internal server error")
	ErrInvalidToken         = errors.New("provided token is invalid")
	ErrTokenExpired         = errors.New("token has expired")
	ErrTokenAlreadyUsed     = errors.New("provided token have already been used")
	ErrInvalidGUID          = errors.New("provided GUID is invalid")
	ErrTokenNotFound        = errors.New("provided token does not exists")
	ErrBadRequest           = errors.New("provided request data is not valid")
)
