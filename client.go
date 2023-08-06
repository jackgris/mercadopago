package mercadopago

import (
	"encoding/json"
	"net/http"
	"time"
)

const (
	baseURLV1 = "https://api.mercadopago.com/v1/"
)

// Client is the API client
type Client struct {
	token      string
	BaseURL    string
	HTTPClient *http.Client
}

// NewClient create a new Client for API interaction
func NewClient(token string) *Client {
	return &Client{
		BaseURL: baseURLV1,
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
