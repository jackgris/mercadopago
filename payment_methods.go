package mercadopago

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

type PaymentMethods []PaymentMethod

type PaymentMethod struct {
	ID                    string        `json:"id"`
	Name                  string        `json:"name"`
	PaymentTypeID         string        `json:"payment_type_id"`
	Status                string        `json:"status"`
	SecureThumbnail       string        `json:"secure_thumbnail"`
	Thumbnail             string        `json:"thumbnail"`
	DeferredCapture       string        `json:"deferred_capture"`
	Settings              []Settings    `json:"settings"`
	AdditionalInfoNeeded  []string      `json:"additional_info_needed"`
	MinAllowedAmount      float64       `json:"min_allowed_amount"`
	MaxAllowedAmount      int           `json:"max_allowed_amount"`
	AccreditationTime     int           `json:"accreditation_time"`
	FinancialInstitutions []interface{} `json:"financial_institutions"`
	ProcessingModes       []string      `json:"processing_modes"`
}

type Settings struct {
	CardNumber   CardNumber   `json:"card_number"`
	Bin          Bin          `json:"bin"`
	SecurityCode SecurityCode `json:"security_code"`
}

type Bin struct {
	Pattern             string `json:"pattern"`
	ExclusionPattern    string `json:"exclusion_pattern"`
	InstallmentsPattern string `json:"installments_pattern"`
}

type CardNumber struct {
	Length     int    `json:"length"`
	Validation string `json:"validation"`
}

type SecurityCode struct {
	Mode         string `json:"mode"`
	Length       int    `json:"length"`
	CardLocation string `json:"card_location"`
}

// PaymentMethods Access to Payment Methods
func (c *Client) PaymentMethods(ctx context.Context) (PaymentMethods, error) {

	url := fmt.Sprintf("%spayment_methods", c.BaseURL)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.Header.Set("Accept", "application/json; charset=utf-8")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.token))

	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	if res.StatusCode < http.StatusOK || res.StatusCode >= http.StatusBadRequest {
		var errRes ErrorResponse
		if err = json.NewDecoder(res.Body).Decode(&errRes); err == nil {
			return nil, &errRes
		}

		return nil, fmt.Errorf("unknown error, status code: %d", res.StatusCode)
	}

	var paymentMethods PaymentMethods
	if err := json.NewDecoder(res.Body).Decode(&paymentMethods); err != nil {
		return nil, errors.New("Can't parse response")
	}

	return paymentMethods, nil
}
