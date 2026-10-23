package repositories

import (
	orderRepo "order-services/repositories/order"
	orderFieldRepo "order-services/repositories/orderfield"
	orderHisotryRepo "order-services/repositories/orderhistory"

	"gorm.io/gorm"
)

type Registry struct {
	db *gorm.DB
}

// GetOrder implements [IRepositoryRegistry].
func (r *Registry) GetOrder() orderRepo.IOrderRepository {
	return orderRepo.NewOrderRepository(r.db)
}

// GetOrderField implements [IRepositoryRegistry].
func (r *Registry) GetOrderField() orderFieldRepo.IOrderFieldRepository {
	return orderFieldRepo.NewOrderFieldRepository(r.db)
}

// GetOrderHistory implements [IRepositoryRegistry].
func (r *Registry) GetOrderHistory() orderHisotryRepo.IOrderHistoryRepository {
	return orderHisotryRepo.NewOrderHistoryRepository(r.db)
}

// GetTx implements [IRepositoryRegistry].
func (r *Registry) GetTx() *gorm.DB {
	return r.db
}

type IRepositoryRegistry interface {
	GetOrder() orderRepo.IOrderRepository
	GetOrderField() orderFieldRepo.IOrderFieldRepository
	GetOrderHistory() orderHisotryRepo.IOrderHistoryRepository
	GetTx() *gorm.DB
}

func NewRepositoryRegistry(db *gorm.DB) IRepositoryRegistry {
	return &Registry{db: db}
}
