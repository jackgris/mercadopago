package mercadopago

import (
	"encoding/json"
	"net/http"
	"time"
)

const (
	BaseURL = "https://api.mercadopago.com/"
)

// Client is the API client
type Client struct {
	token      string
	BaseURL    string
	HTTPClient *http.Client
}

// NewClient create a new Client for API interaction
func NewClient(baseURL, token string) *Client {
	return &Client{
		BaseURL: baseURL,
		token:   token,
		HTTPClient: &http.Client{
			Timeout: time.Minute,
		},
	}
}

// PrettyStruct print out JSON response in a pretty way
func PrettyStruct(data interface{}) (string, error) {
	val, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		return "", err
	}
	return string(val), nil
}
