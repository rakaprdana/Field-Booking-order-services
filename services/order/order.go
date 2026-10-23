package services

import (
	"context"
	"fmt"
	"order-services/clients"
	clientField "order-services/clients/field"
	clientPayment "order-services/clients/payment"
	clientUser "order-services/clients/user"
	"order-services/common/utils"
	"order-services/constants"
	orderErr "order-services/constants/errors/order"
	"order-services/domain/dto"
	"order-services/domain/models"
	"order-services/repositories"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OrderService struct {
	repository repositories.IRepositoryRegistry
	client     clients.IClientRegistry
}

// Create implements [IOrderService].
func (o *OrderService) Create(ctx context.Context, req *dto.OrderRequest) (*dto.OrderResponse, error) {
	var (
		order              *models.Order
		txErr, err         error
		user               = ctx.Value(constants.User).(*clientUser.UserData)
		field              *clientField.FieldData
		paymentResponse    *clientPayment.PaymentData
		orderFieldSchedule = make([]models.OrderField, 0, len(req.FieldScheduleIDs))
		totalAmount        float64
	)

	for _, fieldID := range req.FieldScheduleIDs {
		uuidParse := uuid.MustParse(fieldID)
		field, err = o.client.GetField().GetFieldByUUID(ctx, uuidParse)
		if err != nil {
			return nil, err
		}

		totalAmount += field.PricePerHour
		if field.Status == constants.BookedFieldStatus.String() {
			return nil, orderErr.ErrFieldAlreadyBooked
		}
	}

	err = o.repository.GetTx().Transaction(func(tx *gorm.DB) error {
		order, txErr = o.repository.GetOrder().Create(ctx, tx, &models.Order{
			UserID: user.UUID,
			Amount: totalAmount,
			Date:   time.Now(),
			Status: constants.Pending,
			IsPaid: false,
		})

		if txErr != nil {
			return txErr
		}

		for _, fieldID := range req.FieldScheduleIDs {
			uuidParse := uuid.MustParse(fieldID)
			orderFieldSchedule = append(orderFieldSchedule, models.OrderField{
				OrderID:         order.ID,
				FieldScheduleID: uuidParse,
			})
		}

		txErr = o.repository.GetOrderField().Create(ctx, tx, orderFieldSchedule)
		if txErr != nil {
			return txErr
		}

		txErr = o.repository.GetOrderHistory().Create(ctx, tx, &dto.OrderHistoryRequest{
			Status:  constants.Pending.GetStatusString(),
			OrderID: order.ID,
		})
		if txErr != nil {
			return txErr
		}

		expiredAt := time.Now().Add(1 * time.Hour)
		description := fmt.Sprintf("Pembayaran Sewa %s", field.FieldName)
		paymentResponse, txErr = o.client.GetPayment().CreatePaymentLink(ctx, &dto.PaymentRequest{
			OrderID:     order.UUID,
			ExpiredAt:   expiredAt,
			Amount:      totalAmount,
			Description: description,
			CustomerDetail: &dto.CustomerDetail{
				Name:  user.Name,
				Email: user.Email,
				Phone: user.PhoneNumber,
			},
			ItemDetails: []dto.ItemDetail{
				{
					ID:       uuid.New(),
					Name:     description,
					Amount:   totalAmount,
					Quantity: 1,
				},
			},
		})
		if txErr != nil {
			return txErr
		}

		txErr = o.repository.GetOrder().Update(ctx, tx, &models.Order{
			PaymentID: paymentResponse.UUID,
		}, order.UUID)
		if txErr != nil {
			return txErr
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	response := dto.OrderResponse{
		UUID:        order.UUID,
		Code:        order.Code,
		UserName:    user.Name,
		Amount:      order.Amount,
		Status:      order.Status.GetStatusString(),
		OrderDate:   order.Date,
		PaymentLink: paymentResponse.PaymentLink,
		CreatedAt:   *order.CreatedAt,
		UpdatedAt:   *order.UpdatedAt,
	}

	return &response, nil
}

// GetAllWithPagination implements [IOrderService].
func (o *OrderService) GetAllWithPagination(ctx context.Context, param *dto.OrderRequestParam) (*utils.PaginationResult, error) {
	orders, total, err := o.repository.GetOrder().FindAllWithPagination(ctx, param)
	if err != nil {
		return nil, err
	}

	orderResult := make([]dto.OrderResponse, 0, len(orders))
	for _, order := range orders {
		user, err := o.client.GetUser().GetUserByUUID(ctx, order.UserID)
		if err != nil {
			return nil, err
		}
		orderResult = append(orderResult, dto.OrderResponse{
			UUID:      order.UUID,
			Code:      order.Code,
			UserName:  user.Name,
			Amount:    order.Amount,
			Status:    order.Status.GetStatusString(),
			OrderDate: order.Date,
			CreatedAt: *order.CreatedAt,
			UpdatedAt: *order.UpdatedAt,
		})
	}

	paginationParam := &utils.PaginationParams{
		Count: int(total),
		Page:  param.Page,
		Limit: param.Limit,
		Data:  orderResult,
	}

	response := utils.GeneratePagination(*paginationParam)
	return &response, nil
}

// GetByID implements [IOrderService].
func (o *OrderService) GetByUUID(ctx context.Context, uuid string) (*dto.OrderResponse, error) {
	var (
		order *models.Order
		user  *clientUser.UserData
		err   error
	)

	order, err = o.repository.GetOrder().FindByUUID(ctx, uuid)
	if err != nil {
		return nil, err
	}

	user, err = o.client.GetUser().GetUserByUUID(ctx, order.UserID)
	if err != nil {
		return nil, err
	}

	response := &dto.OrderResponse{
		UUID:      order.UUID,
		Code:      order.Code,
		UserName:  user.Name,
		Amount:    order.Amount,
		Status:    order.Status.GetStatusString(),
		OrderDate: order.Date,
		CreatedAt: *order.CreatedAt,
		UpdatedAt: *order.UpdatedAt,
	}

	return response, nil
}

// GetOrderByUserID implements [IOrderService].
func (o *OrderService) GetOrderByUserID(ctx context.Context) ([]dto.OrderByUserIDResponse, error) {
	var (
		order []models.Order
		user  = ctx.Value(constants.User).(*clientUser.UserData)
		err   error
	)

	order, err = o.repository.GetOrder().FindByUserID(ctx, user.UUID.String())
	if err != nil {
		return nil, err
	}

	orderLists := make([]dto.OrderByUserIDResponse, 0, len(order))
	for _, item := range order {
		payment, err := o.client.GetPayment().GetPaymentByUUID(ctx, item.PaymentID)
		if err != nil {
			return nil, err
		}

		orderLists = append(orderLists, dto.OrderByUserIDResponse{
			Code:        item.Code,
			Amount:      fmt.Sprintf("%s", utils.RupiahFormarter(&item.Amount)),
			Status:      item.Status.GetStatusString(),
			OrderDate:   item.Date.String(),
			PaymentLink: payment.PaymentLink,
			InvoiceLink: payment.InvoiceLink,
		})
	}

	return orderLists, nil
}

// HandlePayment implements [IOrderService].
func (o *OrderService) HandlePayment(ctx context.Context, req *dto.PaymentData) error {
	var (
		err, txErr          error
		order               *models.Order
		orderFieldSchedules []models.OrderField
	)

	status, body := o.mapPaymentStatusToOrder(req)
	err = o.repository.GetTx().Transaction(func(tx *gorm.DB) error {
		txErr = o.repository.GetOrder().Update(ctx, tx, body, req.OrderID)
		if txErr != nil {
			return txErr
		}

		order, txErr = o.repository.GetOrder().FindByUUID(ctx, req.OrderID.String())
		if txErr != nil {
			return txErr
		}

		txErr = o.repository.GetOrderHistory().Create(ctx, tx, &dto.OrderHistoryRequest{
			Status:  status.GetStatusString(),
			OrderID: order.ID,
		})

		if req.Status == constants.SettlementPaymentStatus {
			orderFieldSchedules, txErr = o.repository.GetOrderField().FindByOrderID(ctx, order.ID)
			if txErr != nil {
				return txErr
			}

			fieldScheduleIDs := make([]string, 0, len(orderFieldSchedules))
			for _, item := range orderFieldSchedules {
				fieldScheduleIDs = append(fieldScheduleIDs, item.FieldScheduleID.String())
			}

			txErr = o.client.GetField().UpdateStatus(&dto.UpdateFieldScheduleStatusRequest{
				FieldScheduleIDs: fieldScheduleIDs,
			})
			if txErr != nil {
				return txErr
			}
		}
		return nil
	})

	if err != nil {
		return err
	}
	return nil
}

type IOrderService interface {
	GetAllWithPagination(context.Context, *dto.OrderRequestParam) (*utils.PaginationResult, error)
	GetByUUID(context.Context, string) (*dto.OrderResponse, error)
	GetOrderByUserID(context.Context) ([]dto.OrderByUserIDResponse, error)
	Create(context.Context, *dto.OrderRequest) (*dto.OrderResponse, error)
	HandlePayment(context.Context, *dto.PaymentData) error
}

func NewOrderService(
	repository repositories.IRepositoryRegistry,
	client clients.IClientRegistry,
) IOrderService {
	return &OrderService{
		repository: repository,
		client:     client,
	}
}

func (o *OrderService) mapPaymentStatusToOrder(
	req *dto.PaymentData,
) (constants.OrderStatus, *models.Order) {
	var (
		status constants.OrderStatus
		order  *models.Order
	)

	switch req.Status {
	case constants.SettlementPaymentStatus:
		status = constants.PaymentSuccess
		order = &models.Order{
			IsPaid:    true,
			PaymentID: req.PaymentID,
			PaidAt:    req.PaidAt,
			Status:    status,
		}
	case constants.ExpirePaymentStatus:
		status = constants.Expire
		order = &models.Order{
			IsPaid:    false,
			PaymentID: req.PaymentID,
			PaidAt:    req.PaidAt,
			Status:    status,
		}
	case constants.PendingPaymentStatus:
		status = constants.PendingPayment
		order = &models.Order{
			IsPaid:    false,
			PaymentID: req.PaymentID,
			PaidAt:    req.PaidAt,
			Status:    status,
		}
	}
	return status, order
}
