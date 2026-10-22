package dto

import "order-services/constants"

type PaymentHistoryRequest struct {
	PaymentID uint                          `json:"paymentID"`
	Status    constants.PaymentStatusString `json:"status"`
}
