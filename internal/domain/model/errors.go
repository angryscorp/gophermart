package model

type Error string

func (e Error) Error() string { return string(e) }

const (
	ErrWrongRequestFormat     Error = "wrong request format"
	ErrUnknownInternalError   Error = "unknown internal error"
	ErrUserIsNotAuthenticated Error = "user is not authenticated"
)
