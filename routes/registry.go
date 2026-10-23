package routes

import (
	"order-services/clients"
	controllers "order-services/controllers/http"
	routes "order-services/routes/order"

	"github.com/gin-gonic/gin"
)

type Registry struct {
	controller controllers.IControllerRegistry
	client     clients.IClientRegistry
	group      *gin.RouterGroup
}

// Serve implements [IRouteRegistry].
func (r *Registry) Serve() {
	r.orderRoute().Run()
}

type IRouteRegistry interface {
	Serve()
}

func NewRouteRegistry(
	controller controllers.IControllerRegistry,
	client clients.IClientRegistry,
	group *gin.RouterGroup,
) IRouteRegistry {
	return &Registry{
		controller: controller,
		client:     client,
		group:      group,
	}
}

func (r *Registry) orderRoute() routes.IOrderRoute {
	return routes.NewOrderRoute(r.controller, r.client, r.group)
}
