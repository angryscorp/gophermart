package usecase

type Auth interface {
	SignUp() error
	SignIn() error
}
