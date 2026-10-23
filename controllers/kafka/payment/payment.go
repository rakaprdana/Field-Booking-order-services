package kafka

import (
	"context"
	"encoding/json"
	"order-services/common/utils"
	"order-services/domain/dto"
	"order-services/services"

	"github.com/IBM/sarama"
	"github.com/sirupsen/logrus"
)

const PaymentTopic = "payment-service-callback"

type PaymentKafka struct {
	service services.IServiceRegistry
}

// HandlePayment implements [IPaymentKafka].
func (p *PaymentKafka) HandlePayment(
	ctx context.Context,
	msg *sarama.ConsumerMessage,
) error {
	defer utils.Recover()
	var body dto.PaymentContent
	err := json.Unmarshal(msg.Value, &body)
	if err != nil {
		logrus.Errorf("Failed to unmarshal message: %v", err)
		return err
	}

	data := body.Body.Data
	err = p.service.GetOrder().HandlePayment(ctx, &data)
	if err != nil {
		logrus.Errorf("Failed to handle payment: %v", err)
		return err
	}
	logrus.Infof("Success handle payment")
	return nil
}

type IPaymentKafka interface {
	HandlePayment(context.Context, *sarama.ConsumerMessage) error
}

func NewPaymentKafka(service services.IServiceRegistry) IPaymentKafka {
	return &PaymentKafka{
		service: service,
	}
}
