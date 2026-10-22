package clients

import (
	clientConfig "order-services/clients/config"
	client "order-services/clients/user"
	"order-services/config"
)

type ClientRegistry struct{}

// GetUser implements [IClientRegistry].
func (c *ClientRegistry) GetUser() client.IUserClient {
	return client.NewUserClient(
		clientConfig.NewClientConfig(
			clientConfig.WithBaseURL(config.Config.InternalService.User.Host),
			clientConfig.WithSignatureKey(config.Config.InternalService.User.SignatureKey),
		),
	)
}

type IClientRegistry interface {
	GetUser() client.IUserClient
}

func NewClientRegistry() IClientRegistry {
	return &ClientRegistry{}
}
