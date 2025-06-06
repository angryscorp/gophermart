package usecase

type Balance interface {
	Balance() error
	Withdraw() error
	AllWithdrawals() error
}
