package dto

import "order-services/constants"

type OrderHistoryRequest struct {
	OrderID uint
	Status  constants.OrderStatusString
}
