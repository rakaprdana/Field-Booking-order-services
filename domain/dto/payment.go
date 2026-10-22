package dto

import (
	"time"

	"github.com/google/uuid"
)

type PaymentRequest struct {
	PaymentLink    string          `json:"paymentLink"`
	OrderID        uuid.UUID       `json:"orderID"`
	ExpiredAt      time.Time       `json:"expiredAt"`
	Amount         float64         `json:"amount"`
	Description    string          `json:"description"`
	CustomerDetail *CustomerDetail `json:"customerDetail"`
	ItemDetail     []ItemDetail    `json:"itemDetail"`
}

type CustomerDetail struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Phone string `json:"phone"`
}

type ItemDetail struct {
	ID       string  `json:"id"`
	Amount   float64 `json:"amount"`
	Name     string  `json:"name"`
	Quantity int     `json:"quantity"`
}
