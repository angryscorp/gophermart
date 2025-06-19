package model

type Error string

func (e Error) Error() string { return string(e) }

const (
	ErrUnknownInternalError Error = "unknown internal error"
	ErrUserIsAlreadyExist   Error = "user is already exist"
	ErrUnsufficientFunds    Error = "unsufficient funds"
)
