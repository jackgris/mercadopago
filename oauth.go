package mercadopago

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

type RequestAccessToken struct {
	ClientSecret string `json:"client_secret"`
	ClientID     string `json:"client_id"`
	GrantType    string `json:"grant_type"`
}

type AccessToken struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
	Scope       string `json:"scope"`
	UserID      int    `json:"user_id"`
}

// GetAccessToken return our access token to start operating with the endpoints.
// For obtain the client secret and client id you need to go to your integrations: https://www.mercadopago.com.ar/developers/panel/app
// Choose one of then and in the production credentials you will find them. (You can't use the test credentials for this)
func (c *Client) GetAccessToken(ctx context.Context, clientId, clientSecret string) (*AccessToken, error) {
	data := RequestAccessToken{
		ClientSecret: clientSecret,
		ClientID:     clientId,
		GrantType:    "client_credentials",
	}

	body, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	url := fmt.Sprintf("%soauth/token", c.BaseURL[:len(c.BaseURL)-3])

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(body))
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

	var accessToken AccessToken
	if err := json.NewDecoder(res.Body).Decode(&accessToken); err != nil {
		return nil, errors.New("Can't parse response")
	}

	return &accessToken, nil
}
