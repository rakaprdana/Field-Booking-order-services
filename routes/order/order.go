package routes

import (
	"order-services/clients"
	"order-services/constants"
	controllers "order-services/controllers/http"
	"order-services/middlewares"

	"github.com/gin-gonic/gin"
)

type OrderRoute struct {
	controller controllers.IControllerRegistry
	client     clients.IClientRegistry
	group      *gin.RouterGroup
}

// Run implements [IOrderRoute].
func (o *OrderRoute) Run() {
	group := o.group.Group("/order")
	group.Use(middlewares.Authenticated())
	group.GET("", middlewares.CheckRole([]string{
		constants.Admin,
		constants.Customer,
	}, o.client), o.controller.GetOrder().GetAllWithPagination)
	group.GET("/:uuid", middlewares.CheckRole([]string{
		constants.Admin,
		constants.Customer,
	}, o.client), o.controller.GetOrder().GetByUUID)
	group.GET("/user", middlewares.CheckRole([]string{
		constants.Customer,
	}, o.client), o.controller.GetOrder().GetOrderByUserID)
	group.POST("", middlewares.CheckRole([]string{
		constants.Customer,
	}, o.client), o.controller.GetOrder().Create)
}

type IOrderRoute interface {
	Run()
}

func NewOrderRoute(
	controller controllers.IControllerRegistry,
	client clients.IClientRegistry,
	group *gin.RouterGroup,
) IOrderRoute {
	return &OrderRoute{
		controller: controller,
		client:     client,
		group:      group,
	}
}
