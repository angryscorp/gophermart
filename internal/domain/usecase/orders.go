package usecase

type Orders interface {
	Upload() error
	AllOrders() error
}
