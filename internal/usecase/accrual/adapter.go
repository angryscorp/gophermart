package accrual

import (
	"encoding/json"
	"fmt"
	"github.com/angryscorp/gophermart/internal/domain/model"
	"github.com/angryscorp/gophermart/internal/domain/usecase"
	"github.com/rs/zerolog"
	"net/http"
	"strings"
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
	if !strings.HasPrefix(accrualAddress, "http://") && !strings.HasPrefix(accrualAddress, "https://") {
		if strings.HasPrefix(accrualAddress, ":") {
			accrualAddress = "http://localhost" + accrualAddress
		} else {
			accrualAddress = "http://" + accrualAddress
		}
	}

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
	defer func() {
		if err := resp.Body.Close(); err != nil {
			a.logger.Error().Err(err).Msg("failed to close response body")
		}
	}()

	if resp.StatusCode != http.StatusOK {
		a.logger.Error().Int("status", resp.StatusCode).Str("order", orderNumber).Msg("non-200 response from accrual service")
		return model.Accrual{}, fmt.Errorf("accrual service returned status %d", resp.StatusCode)
	}

	var accrualResp model.Accrual
	if err := json.NewDecoder(resp.Body).Decode(&accrualResp); err != nil {
		a.logger.Error().Err(err).Str("order", orderNumber).Msg("failed to decode accrual response")
		return model.Accrual{}, err
	}

	accrualValue := 0
	if accrualResp.Accrual != nil {
		accrualValue = *accrualResp.Accrual
	}

	a.logger.Info().Str("order", orderNumber).Str("status", string(accrualResp.Status)).Int("accrual", accrualValue).Msg("accrual info received")

	return accrualResp, nil
}
