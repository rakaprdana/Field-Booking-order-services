package repositories

import (
	"context"
	errWrap "order-services/common/errors"
	constErr "order-services/constants/errors"
	"order-services/domain/dto"
	"order-services/domain/models"

	"gorm.io/gorm"
)

type OrderHistoryRepository struct {
	db *gorm.DB
}

// Create implements [IOrderHistoryRepository].
func (o *OrderHistoryRepository) Create(
	ctx context.Context,
	tx *gorm.DB,
	param *dto.OrderHistoryRequest,
) error {
	orderHistory := &models.OrderHistory{
		OrderID: param.OrderID,
		Status:  param.Status,
	}

	err := tx.
		WithContext(ctx).
		Create(&orderHistory).
		Error

	if err != nil {
		return errWrap.WrapError(constErr.ErrSQLError)
	}

	return nil
}

type IOrderHistoryRepository interface {
	Create(context.Context, *gorm.DB, *dto.OrderHistoryRequest) error
}

func NewOrderHistoryRepository(db *gorm.DB) IOrderHistoryRepository {
	return &OrderHistoryRepository{db: db}
}
