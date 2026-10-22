package clients

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	config "order-services/clients/config"
	"order-services/common/utils"
	configApp "order-services/config"
	"order-services/constants"
	"order-services/domain/dto"
	"time"

	"github.com/google/uuid"
)

type PaymentClient struct {
	client config.IClientConfig
}

// CreatePaymentLinkt implements [IPaymentClient].
func (p *PaymentClient) CreatePaymentLink(ctx context.Context, req *dto.PaymentRequest) (*PaymentData, error) {
	unixTime := time.Now().Unix()
	generateAPIkey := fmt.Sprintf("%s:%s:%d",
		configApp.Config.AppName,
		p.client.SignatureKey(),
		unixTime,
	)

	apiKey := utils.GenerateSHA256(generateAPIkey)
	token := ctx.Value(constants.Token).(string)
	bearerToken := fmt.Sprintf("Bearer %s", token)

	body, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	resp, bodyResp, errs := p.client.Client().Clone().
		Post(fmt.Sprintf("%s/api/v1/payment/", p.client.BaseURL())).
		Set(constants.Authorization, bearerToken).
		Set(constants.XApiKey, apiKey).
		Set(constants.XServiceName, configApp.Config.AppName).
		Set(constants.XRequestAt, fmt.Sprintf("%d", unixTime)).
		Send(string(body)).
		End()

	if len(errs) > 0 {
		return nil, errs[0]
	}
	var res PaymentResponse
	if resp.StatusCode != http.StatusCreated {
		err = json.Unmarshal([]byte(bodyResp), &res)
		if err != nil {
			return nil, err
		}

		paymentErr := fmt.Errorf("payment response: %s", res.Message)
		return nil, paymentErr
	}

	err = json.Unmarshal([]byte(bodyResp), &res)
	if err != nil {
		return nil, err
	}

	return &res.Data, nil
}

// GetPaymentByUUID implements [IPaymentClient].
func (p *PaymentClient) GetPaymentByUUID(ctx context.Context, uuid uuid.UUID) (*PaymentData, error) {
	unixTime := time.Now().Unix()
	generateAPIkey := fmt.Sprintf("%s:%s:%d",
		configApp.Config.AppName,
		p.client.SignatureKey(),
		unixTime,
	)

	apiKey := utils.GenerateSHA256(generateAPIkey)
	token := ctx.Value(constants.Token).(string)
	bearerToken := fmt.Sprintf("Bearer %s", token)

	var res PaymentResponse
	req := p.client.Client().Clone().
		Set(constants.Authorization, bearerToken).
		Set(constants.XApiKey, apiKey).
		Set(constants.XServiceName, configApp.Config.AppName).
		Set(constants.XRequestAt, fmt.Sprintf("%d", unixTime)).
		Get(fmt.Sprintf("%s/api/v1/payment/%s", p.client.BaseURL(), uuid))

	resp, _, errs := req.EndStruct(&res)
	if len(errs) > 0 {
		return nil, errs[0]
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Payment response :%s", res.Message)
	}

	return &res.Data, nil
}

type IPaymentClient interface {
	GetPaymentByUUID(context.Context, uuid.UUID) (*PaymentData, error)
	CreatePaymentLink(context.Context, *dto.PaymentRequest) (*PaymentData, error)
}

func NewPaymentClient(client config.IClientConfig) IPaymentClient {
	return &PaymentClient{client: client}
}
