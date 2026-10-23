package controllers

import (
	"net/http"
	errValidation "order-services/common/errors"
	response "order-services/common/reponses"
	"order-services/domain/dto"
	"order-services/services"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type OrderController struct {
	service services.IServiceRegistry
}

// Create implements [IOrderController].
func (o *OrderController) Create(c *gin.Context) {
	var (
		req dto.OrderRequest
		ctx = c.Request.Context()
	)

	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.HttpResponse(response.ParamHTTPResp{
			Code: http.StatusBadRequest,
			Err:  err,
			Gin:  c,
		})
		return
	}

	validate := validator.New()
	if err = validate.Struct(req); err != nil {
		errMessage := http.StatusText(http.StatusUnprocessableEntity)
		errResponse := errValidation.ErrValidationResponse(err)
		response.HttpResponse(response.ParamHTTPResp{
			Err:     err,
			Code:    http.StatusUnprocessableEntity,
			Message: &errMessage,
			Data:    errResponse,
			Gin:     c,
		})

		return
	}

	result, err := o.service.GetOrder().Create(ctx, &req)
	if err != nil {
		response.HttpResponse(response.ParamHTTPResp{
			Code: http.StatusBadRequest,
			Err:  err,
			Gin:  c,
		})

		return
	}

	response.HttpResponse(response.ParamHTTPResp{
		Code: http.StatusCreated,
		Data: result,
		Gin:  c,
	})

}

// GetAllWithPagination implements [IOrderController].
func (o *OrderController) GetAllWithPagination(ctx *gin.Context) {
	var param dto.OrderRequestParam
	err := ctx.ShouldBindQuery(param)

	if err != nil {
		response.HttpResponse(response.ParamHTTPResp{
			Code: http.StatusBadRequest,
			Err:  err,
			Gin:  ctx,
		})
		return
	}

	validate := validator.New()
	if err = validate.Struct(param); err != nil {
		errMessage := http.StatusText(http.StatusUnprocessableEntity)
		errResponse := errValidation.ErrValidationResponse(err)
		response.HttpResponse(response.ParamHTTPResp{
			Err:     err,
			Code:    http.StatusUnprocessableEntity,
			Message: &errMessage,
			Data:    errResponse,
			Gin:     ctx,
		})

		return
	}

	result, err := o.service.GetOrder().GetAllWithPagination(ctx, &param)
	if err != nil {
		response.HttpResponse(response.ParamHTTPResp{
			Code: http.StatusBadRequest,
			Err:  err,
			Gin:  ctx,
		})

		return
	}

	response.HttpResponse(response.ParamHTTPResp{
		Code: http.StatusOK,
		Data: result,
		Gin:  ctx,
	})
}

// GetByUUID implements [IOrderController].
func (o *OrderController) GetByUUID(ctx *gin.Context) {
	result, err := o.service.GetOrder().GetByUUID(ctx, ctx.Param("uuid"))
	if err != nil {
		response.HttpResponse(response.ParamHTTPResp{
			Code: http.StatusBadRequest,
			Err:  err,
			Gin:  ctx,
		})
		return
	}

	response.HttpResponse(response.ParamHTTPResp{
		Code: http.StatusOK,
		Data: result,
		Gin:  ctx,
	})
}

// GetOrderByUserID implements [IOrderController].
func (o *OrderController) GetOrderByUserID(ctx *gin.Context) {
	result, err := o.service.GetOrder().GetOrderByUserID(ctx.Request.Context())
	if err != nil {
		response.HttpResponse(response.ParamHTTPResp{
			Code: http.StatusBadRequest,
			Err:  err,
			Gin:  ctx,
		})
		return
	}

	response.HttpResponse(response.ParamHTTPResp{
		Code: http.StatusOK,
		Data: result,
		Gin:  ctx,
	})
}

type IOrderController interface {
	GetAllWithPagination(*gin.Context)
	GetByUUID(*gin.Context)
	GetOrderByUserID(*gin.Context)
	Create(*gin.Context)
}

func NewOrderController(service services.IServiceRegistry) IOrderController {
	return &OrderController{
		service: service,
	}
}
