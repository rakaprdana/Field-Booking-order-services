package repositories

import (
	"context"
	errWrap "order-services/common/errors"
	constErr "order-services/constants/errors"
	"order-services/domain/models"

	"gorm.io/gorm"
)

type OrderFieldRepository struct {
	db *gorm.DB
}

// Create implements [IOrderFieldRepository].
func (o *OrderFieldRepository) Create(ctx context.Context, tx *gorm.DB, param []models.OrderField) error {
	err := tx.
		WithContext(ctx).
		Create(&param).
		Error

	if err != nil {
		return errWrap.WrapError(constErr.ErrSQLError)
	}

	return nil
}

// FindByOrderID implements [IOrderFieldRepository].
func (o *OrderFieldRepository) FindByOrderID(ctx context.Context, orderID uint) ([]models.OrderField, error) {
	var orderFields []models.OrderField
	err := o.db.
		WithContext(ctx).
		Where("order_id = ?", orderID).
		First(&orderFields).
		Error
	if err != nil {
		return nil, errWrap.WrapError(constErr.ErrSQLError)
	}

	return orderFields, nil
}

type IOrderFieldRepository interface {
	FindByOrderID(context.Context, uint) ([]models.OrderField, error)
	Create(context.Context, *gorm.DB, []models.OrderField) error
}

func NewOrderFieldRepository(db *gorm.DB) IOrderFieldRepository {
	return &OrderFieldRepository{db: db}
}
