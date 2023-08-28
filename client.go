package mercadopago

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

const (
	BaseURL = "https://api.mercadopago.com/"
)

// Client is the API client
type Client struct {
	token         string
	BaseURL       string
	HTTPClient    *http.Client
	HTTPTransport transport
}

// NewClient create a new Client for API interaction
func NewClient(rawURL, token string) *Client {
	baseURL, err := url.ParseRequestURI(rawURL)
	if err != nil {
		return nil
	}

	t := transport{
		header:  http.Header{},
		baseUrl: *baseURL,
	}

	client := &http.Client{
		Timeout:   time.Minute,
		Transport: t,
	}

	c := Client{
		BaseURL:       baseURL.String(),
		token:         token,
		HTTPClient:    client,
		HTTPTransport: t,
	}

	c.HTTPTransport.header.Set("Content-Type", "application/json; charset=utf-8")
	c.HTTPTransport.header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	return &c
}

type transport struct {
	header  http.Header
	baseUrl url.URL
}

func (t transport) RoundTrip(request *http.Request) (*http.Response, error) {
	for headerName, values := range t.header {
		for _, val := range values {
			request.Header.Add(headerName, val)
		}
	}
	request.URL = t.baseUrl.ResolveReference(request.URL)
	return http.DefaultTransport.RoundTrip(request)
}

// PrettyStruct print out JSON response in a pretty way
func PrettyStruct(data interface{}) (string, error) {
	val, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		return "", err
	}
	return string(val), nil
}
