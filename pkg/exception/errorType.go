package exception

import "errors"

var (
	ErrUnauthorizedRequest      = errors.New("unauthorized request")
	ErrInvalidParameter         = errors.New("invalid parameter")
	ErrIncorrectConfirmPassword = errors.New("incorrect confirm password")
	ErrWrongPassword            = errors.New("wrong password")
	ErrUnableToDeleteUser       = errors.New("unable to delete user")
	ErrUnexpectedSigningMethod  = errors.New("unexpected signing method")
	ErrInvalidTokenClaims       = errors.New("invalid token claims")
	ErrMissingParams            = errors.New("missing params")
)

type ErrValidateResult struct {
	IsError bool
	Field   string
	Tag     string
	Value   interface{}
}

type ErrResponseCtx struct {
	IsError    bool        `json:"isError"`
	StatusCode int         `json:"statusCode"`
	Message    string      `json:"message"`
	Error      interface{} `json:"error"`
}
