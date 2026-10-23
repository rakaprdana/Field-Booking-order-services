package kafka

import (
	kafka "order-services/controllers/kafka/payment"
	"order-services/services"
)

type KafkaRegistry struct {
	service services.IServiceRegistry
}

// GetPayment implements [IKafkaRegistry].
func (k *KafkaRegistry) GetPayment() kafka.IPaymentKafka {
	return kafka.NewPaymentKafka(k.service)
}

type IKafkaRegistry interface {
	GetPayment() kafka.IPaymentKafka
}

func NewKafkaRegistry(service services.IServiceRegistry) IKafkaRegistry {
	return &KafkaRegistry{
		service: service,
	}
}
