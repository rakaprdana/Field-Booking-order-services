package kafka

import (
	"order-services/config"
	"order-services/controllers/kafka"
	kafkaPayment "order-services/controllers/kafka/payment"

	"golang.org/x/exp/slices"
)

type Kafka struct {
	consumer *ConsumerGroup
	kafka    kafka.IKafkaRegistry
}

// Register implements [IKafka].
func (k *Kafka) Register() {
	k.paymentHandler()
}

type IKafka interface {
	Register()
}

func NewKafkaConsumer(
	consumer *ConsumerGroup,
	kafka kafka.IKafkaRegistry,
) IKafka {
	return &Kafka{
		consumer: consumer,
		kafka:    kafka,
	}
}

func (k *Kafka) paymentHandler() {
	if slices.Contains(config.Config.Kafka.Topic, kafkaPayment.PaymentTopic) {
		k.consumer.RegisterHandler(kafkaPayment.PaymentTopic, k.kafka.GetPayment().HandlePayment)
	}
}
