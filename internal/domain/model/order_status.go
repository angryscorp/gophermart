package model

type OrderStatus string

const (
	OrderStatusNew        OrderStatus = "NEW"
	OrderStatusProcessing OrderStatus = "PROCESSING"
	OrderStatusInvalid    OrderStatus = "INVALID"
	OrderStatusProcessed  OrderStatus = "PROCESSED"
)

func NewOrderStatus(status string) OrderStatus {
	switch status {
	case "NEW":
		return OrderStatusNew
	case "PROCESSING":
		return OrderStatusProcessing
	case "INVALID":
		return OrderStatusInvalid
	case "PROCESSED":
		return OrderStatusProcessed
	default:
		return OrderStatusNew
	}
}
