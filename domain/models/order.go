package models

import (
	"order-services/constants"
	"time"

	"github.com/google/uuid"
)

type Order struct {
	ID        uint                  `gorm:"primaryKey;autoIncrement"`
	UUID      uuid.UUID             `gorm:"type:uuid;not null"`
	Code      string                `gorm:"type:varchar(30);not null"`
	UserID    uuid.UUID             `gorm:"type:uuid;not null"`
	PaymentID uuid.UUID             `gorm:"type:uuid;not null"`
	Amount    float64               `gorm:"decimal(10,2);not null"`
	Status    constants.OrderStatus `gorm:"type:int;not null"`
	Date      time.Time             `gorm:"type:timestamp;not null"`
	PaidAt    *time.Time            `gorm:"type:timestamp"`
	IsPaid    bool                  `gorm:"type:boolean;default:false;not null"`
	CreatedAt *time.Time
	UpdatedAt *time.Time
}
