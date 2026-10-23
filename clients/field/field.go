package clients

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	client "order-services/clients/config"
	"order-services/common/utils"
	"order-services/config"
	"order-services/constants"
	"order-services/domain/dto"
	"time"

	"github.com/google/uuid"
)

type FieldClient struct {
	client client.IClientConfig
}

// GetFieldByUUID implements [IFieldClient].
func (f *FieldClient) GetFieldByUUID(ctx context.Context, uuid uuid.UUID) (*FieldData, error) {
	unixTime := time.Now().Unix()
	generateAPIkey := fmt.Sprintf("%s:%s:%d",
		config.Config.AppName,
		f.client.SignatureKey(),
		unixTime,
	)

	apiKey := utils.GenerateSHA256(generateAPIkey)
	token := ctx.Value(constants.Token).(string)
	bearerToken := fmt.Sprintf("Bearer %s", token)

	var res FieldResponse
	req := f.client.Client().Clone().
		Set(constants.Authorization, bearerToken).
		Set(constants.XApiKey, apiKey).
		Set(constants.XServiceName, config.Config.AppName).
		Set(constants.XRequestAt, fmt.Sprintf("%d", unixTime)).
		Get(fmt.Sprintf("%s/api/v1/field/schedule/%s", f.client.BaseURL(), uuid))

	resp, _, errs := req.EndStruct(&res)
	if len(errs) > 0 {
		return nil, errs[0]
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("User response :%s", res.Message)
	}

	return &res.Data, nil
}

// UpdateStatus implements [IFieldClient].
func (f *FieldClient) UpdateStatus(req *dto.UpdateFieldScheduleStatusRequest) error {
	unixTime := time.Now().Unix()
	generateAPIkey := fmt.Sprintf("%s:%s:%d",
		config.Config.AppName,
		f.client.SignatureKey(),
		unixTime,
	)

	apiKey := utils.GenerateSHA256(generateAPIkey)

	body, err := json.Marshal(req)
	if err != nil {
		return err
	}

	resp, bodyResp, errs := f.client.Client().Clone().
		Post(fmt.Sprintf("%s/api/v1/field/schedule/status", f.client.BaseURL())).
		Set(constants.XApiKey, apiKey).
		Set(constants.XServiceName, config.Config.AppName).
		Set(constants.XRequestAt, fmt.Sprintf("%d", unixTime)).
		Send(string(body)).
		End()

	if len(errs) > 0 {
		return errs[0]
	}
	var res FieldResponse
	if resp.StatusCode != http.StatusCreated {
		err = json.Unmarshal([]byte(bodyResp), &res)
		if err != nil {
			return err
		}

		fieldErr := fmt.Errorf("field response: %s", res.Message)
		return fieldErr
	}

	err = json.Unmarshal([]byte(bodyResp), &res)
	if err != nil {
		return err
	}

	return nil
}

type IFieldClient interface {
	GetFieldByUUID(context.Context, uuid.UUID) (*FieldData, error)
	UpdateStatus(*dto.UpdateFieldScheduleStatusRequest) error
}

func NewFieldClient(client client.IClientConfig) IFieldClient {
	return &FieldClient{client: client}
}
