package services

import (
	"order-services/clients"
	"order-services/repositories"
	services "order-services/services/order"
)

type Registry struct {
	repository repositories.IRepositoryRegistry
	client     clients.IClientRegistry
}

// GetOrder implements [IServiceRegistry].
func (r *Registry) GetOrder() services.IOrderService {
	return services.NewOrderService(r.repository, r.client)
}

type IServiceRegistry interface {
	GetOrder() services.IOrderService
}

func NewRegistry(repository repositories.IRepositoryRegistry, client clients.IClientRegistry) IServiceRegistry {
	return &Registry{
		repository: repository,
		client:     client,
	}
}
