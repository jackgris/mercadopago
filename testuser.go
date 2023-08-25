package mercadopago

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

type ResquestTestUser struct {
	SiteID      string `json:"site_id"`
	Description string `json:"description"`
}

type TestUser struct {
	ID              int    `json:"id"`
	Nickname        string `json:"nickname"`
	Password        string `json:"password"`
	SiteStatus      string `json:"site_status"`
	SiteID          string `json:"site_id"`
	Description     string `json:"description"`
	Email           string `json:"email"`
	DateCreated     string `json:"date_created"`
	DateLastUpdated string `json:"date_last_updated"`
}

// GetTestUser use the endpoint that handles http requests to create a test user. And return data of a user.
// We can use that data like email and password to interact with others endpoint of MercadoPago.
// Possibles ID of the site where the test user will be created:
// MPE: Mercado Libre Perú
// MLU: Mercado Libre Uruguay
// MLA: Mercado Libre Argentina
// MLC: Mercado Libre Chile
// MCO: Mercado Libre Colombia
// MLB: Mercado Libre Brasil
// MLM: Mercado Libre México
func (c *Client) GetTestUser(ctx context.Context, accessToken, siteId, description string) (*TestUser, error) {
	data := ResquestTestUser{
		SiteID:      siteId,
		Description: description,
	}

	body, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	url := fmt.Sprintf("%susers/test_user", c.BaseURL)

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))

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

	var testUser TestUser
	if err := json.NewDecoder(res.Body).Decode(&testUser); err != nil {
		return nil, errors.New("Can't parse response")
	}

	return &testUser, nil
}
