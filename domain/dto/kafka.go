package dto

import (
	"time"

	"github.com/google/uuid"
)

type KafkaEvent struct {
	Name string `json:"name"`
}

type KafkaMetaData struct {
	Sender    string `json:"sender"`
	SendingAt string `json:"senderAt"`
}

type KafkaData struct {
	OrderID   uuid.UUID  `json:"orderID"`
	PaymentID uuid.UUID  `json:"paymentID"`
	Status    string     `json:"status"`
	ExpiredAt time.Time  `json:"expiredAt"`
	PaidAt    *time.Time `json:"paidAt"`
}
type DataType string

type KafkaBody[T any] struct {
	Type DataType `json:"type"`
	Data T        `json:"data"`
}

type KafkaMessage[T any] struct {
	Event    KafkaEvent    `json:"event"`
	Metadata KafkaMetaData `json:"metadata"`
	Body     KafkaBody[T]  `json:"body"`
}
