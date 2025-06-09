package accrual

import (
	"encoding/json"
	"github.com/angryscorp/gophermart/internal/domain/model"
	"github.com/angryscorp/gophermart/internal/domain/usecase"
	"github.com/rs/zerolog"
	"io"
	"net/http"
)

type Adapter struct {
	client         *http.Client
	logger         zerolog.Logger
	accrualAddress string
}

var _ usecase.Accrual = (*Adapter)(nil)

func NewAdapter(
	client *http.Client,
	logger zerolog.Logger,
	accrualAddress string,
) Adapter {
	return Adapter{
		client:         client,
		logger:         logger,
		accrualAddress: accrualAddress,
	}
}

func (a Adapter) Status(orderNumber string) (model.Accrual, error) {
	a.logger.Info().Str("order", orderNumber).Msg("fetching accrual info")

	resp, err := a.client.Get(a.accrualAddress + "/api/orders/" + orderNumber)
	if err != nil {
		a.logger.Error().Err(err).Str("order", orderNumber).Msg("failed to get accrual info")
		return model.Accrual{}, err
	}
	defer func(body io.ReadCloser) {
		_ = body.Close()
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		a.logger.Error().Int("status", resp.StatusCode).Str("order", orderNumber).Msg("non-200 response from accrual service")
		return model.Accrual{}, err
	}

	var accrualResp model.Accrual
	if err := json.NewDecoder(resp.Body).Decode(&accrualResp); err != nil {
		a.logger.Error().Err(err).Str("order", orderNumber).Msg("failed to decode accrual response")
		return model.Accrual{}, err
	}

	a.logger.Info().Str("order", orderNumber).Str("status", string(accrualResp.Status)).Int("accrual", accrualResp.Accrual).Msg("accrual info received")

	return accrualResp, nil
}
