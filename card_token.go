package mercadopago

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

type RequestCardToken struct {
	CardNumber      string `json:"card_number"`
	ExpirationMonth int    `json:"expiration_month"`
	ExpirationYear  int    `json:"expiration_year"`
	SecurityCode    string `json:"security_code"`
	Cardholder      struct {
		Name string `json:"name"`
	} `json:"cardholder"`
}

type CardToken struct {
	ID              string `json:"id"`
	FirstSixDigits  string `json:"first_six_digits"`
	ExpirationMonth int    `json:"expiration_month"`
	ExpirationYear  int    `json:"expiration_year"`
	LastFourDigits  string `json:"last_four_digits"`
	Cardholder      struct {
		Identification struct {
		} `json:"identification"`
	} `json:"cardholder"`
	Status             string `json:"status"`
	DateCreated        string `json:"date_created"`
	DateLastUpdated    string `json:"date_last_updated"`
	DateDue            string `json:"date_due"`
	LuhnValidation     bool   `json:"luhn_validation"`
	LiveMode           bool   `json:"live_mode"`
	RequireEsc         bool   `json:"require_esc"`
	CardNumberLength   int    `json:"card_number_length"`
	SecurityCodeLength int    `json:"security_code_length"`
}

// GetCardToken will retrieve all the data from the credit card, including the ID necessary to make the payments.
func (c *Client) GetCardToken(ctx context.Context, data RequestCardToken) (*CardToken, error) {

	body, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	url := fmt.Sprintf("%sv1/card_tokens", c.BaseURL)

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
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

	var cardToken CardToken
	if err := json.NewDecoder(res.Body).Decode(&cardToken); err != nil {
		return nil, errors.New("Can't parse response")
	}

	return &cardToken, nil
}
