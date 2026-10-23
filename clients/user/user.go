package client

import (
	"context"
	"fmt"
	"net/http"
	client "order-services/clients/config"
	"order-services/common/utils"
	"order-services/config"
	"order-services/constants"
	"time"

	"github.com/google/uuid"
)

type UserClient struct {
	client client.IClientConfig
}

type IUserClient interface {
	GetUserByToken(context.Context) (*UserData, error)
	GetUserByUUID(context.Context, uuid.UUID) (*UserData, error)
}

func NewUserClient(client client.IClientConfig) IUserClient {
	return &UserClient{client: client}
}

// GetUserByToken implements [IUserClient].
func (u *UserClient) GetUserByToken(ctx context.Context) (*UserData, error) {
	unixTime := time.Now().Unix()
	generateAPIkey := fmt.Sprintf("%s:%s:%d",
		config.Config.AppName,
		u.client.SignatureKey(),
		unixTime,
	)

	apiKey := utils.GenerateSHA256(generateAPIkey)
	token := ctx.Value(constants.Token).(string)
	bearerToken := fmt.Sprintf("Bearer %s", token)

	var res UserResponse
	req := u.client.Client().Clone().
		Set(constants.Authorization, bearerToken).
		Set(constants.XApiKey, apiKey).
		Set(constants.XServiceName, config.Config.AppName).
		Set(constants.XRequestAt, fmt.Sprintf("%d", unixTime)).
		Get(fmt.Sprintf("%s/api/v1/auth/user", u.client.BaseURL()))

	resp, _, errs := req.EndStruct(&res)
	if len(errs) > 0 {
		return nil, errs[0]
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("User response :%s", res.Message)
	}

	return &res.Data, nil
}

// GetUserByUUID implements [IUserClient].
func (u *UserClient) GetUserByUUID(ctx context.Context, uuid uuid.UUID) (*UserData, error) {
	unixTime := time.Now().Unix()
	generateAPIkey := fmt.Sprintf("%s:%s:%d",
		config.Config.AppName,
		u.client.SignatureKey(),
		unixTime,
	)

	apiKey := utils.GenerateSHA256(generateAPIkey)
	token := ctx.Value(constants.Token).(string)
	bearerToken := fmt.Sprintf("Bearer %s", token)

	var res UserResponse
	req := u.client.Client().Clone().
		Set(constants.Authorization, bearerToken).
		Set(constants.XApiKey, apiKey).
		Set(constants.XServiceName, config.Config.AppName).
		Set(constants.XRequestAt, fmt.Sprintf("%d", unixTime)).
		Get(fmt.Sprintf("%s/api/v1/auth/%s", u.client.BaseURL(), uuid))

	resp, _, errs := req.EndStruct(&res)
	if len(errs) > 0 {
		return nil, errs[0]
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("User response :%s", res.Message)
	}

	return &res.Data, nil
}
