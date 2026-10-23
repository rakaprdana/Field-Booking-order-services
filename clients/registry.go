package clients

import (
	clientConfig "order-services/clients/config"
	clientField "order-services/clients/field"
	clientPayment "order-services/clients/payment"
	clientUser "order-services/clients/user"
	"order-services/config"
)

type ClientRegistry struct{}

// GetField implements [IClientRegistry].
func (c *ClientRegistry) GetField() clientField.IFieldClient {
	return clientField.NewFieldClient(
		clientConfig.NewClientConfig(
			clientConfig.WithBaseURL(config.Config.InternalService.Field.Host),
			clientConfig.WithSignatureKey(config.Config.InternalService.Field.SignatureKey),
		),
	)

}

// GetPayment implements [IClientRegistry].
func (c *ClientRegistry) GetPayment() clientPayment.IPaymentClient {
	return clientPayment.NewPaymentClient(
		clientConfig.NewClientConfig(
			clientConfig.WithBaseURL(config.Config.InternalService.Payment.Host),
			clientConfig.WithSignatureKey(config.Config.InternalService.Payment.SignatureKey),
		),
	)
}

// GetUser implements [IClientRegistry].
func (c *ClientRegistry) GetUser() clientUser.IUserClient {
	return clientUser.NewUserClient(
		clientConfig.NewClientConfig(
			clientConfig.WithBaseURL(config.Config.InternalService.User.Host),
			clientConfig.WithSignatureKey(config.Config.InternalService.User.SignatureKey),
		),
	)
}

type IClientRegistry interface {
	GetUser() clientUser.IUserClient
	GetPayment() clientPayment.IPaymentClient
	GetField() clientField.IFieldClient
}

func NewClientRegistry() IClientRegistry {
	return &ClientRegistry{}
}
