package mercadopago_test

import (
	"testing"

	"github.com/jackgris/mercadopago"
)

func TestNewClient(t *testing.T) {
	c := mercadopago.NewClient("hhttpp//asdasd", "")
	if c != nil {
		t.Error("Client should be nil because receive an empty string as base URL.")
	}

	c = mercadopago.NewClient(mercadopago.BaseURL, "")
	if c == nil {
		t.Error("Client should not be nil because receive a right string as base URL.")
	}
}
