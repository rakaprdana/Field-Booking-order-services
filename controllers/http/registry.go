package controllers

import (
	controller "order-services/controllers/http/order"
	"order-services/services"
)

type Registry struct {
	service services.IServiceRegistry
}

// GetOrder implements [IControllerRegistry].
func (r *Registry) GetOrder() controller.IOrderController {
	return controller.NewOrderController(r.service)
}

type IControllerRegistry interface {
	GetOrder() controller.IOrderController
}

func NewContollerRegistry(
	service services.IServiceRegistry,
) IControllerRegistry {
	return &Registry{
		service: service,
	}
}
