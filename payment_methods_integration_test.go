//go:build integration
// +build integration

package mercadopago_test

import (
	"context"
	"errors"
	"net/http"
	"os"
	"testing"

	"github.com/jackgris/mercadopago"
	"github.com/joho/godotenv"
)

func TestIntegrationPaymentMethods(t *testing.T) {

	err := godotenv.Load("test.env")
	if err != nil {
		t.Fatal("Error loading .env file")
	}
	accessToken := os.Getenv("ACCESS_TOKEN")

	tests := []struct {
		name        string
		respStatus  int
		accessToken string
		expectedErr *mercadopago.ErrorResponse
	}{
		{
			name:        "Successful response",
			respStatus:  http.StatusOK,
			accessToken: accessToken,
			expectedErr: nil,
		},
		{
			name:        "Unauthorized",
			respStatus:  http.StatusUnauthorized,
			accessToken: "",
			expectedErr: &mercadopago.ErrorResponse{
				Message: "neither a public key or caller id were provided",
				Errors:  "unauthorized_scopes",
				Status:  http.StatusUnauthorized,
			},
		},
		{
			name:        "Not Found",
			respStatus:  http.StatusUnauthorized,
			accessToken: "TEST-1237334875597470-080318-abc3babd65d6d886dd1193889f2b85a4-4708",
			expectedErr: &mercadopago.ErrorResponse{
				Message: "invalid_token",
				Errors:  "not_found",
				Status:  http.StatusUnauthorized,
			},
		},
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {

			client := mercadopago.NewClient(mercadopago.BaseURL, tt.accessToken)
			got, err := client.PaymentMethods(context.Background())
			if tt.respStatus == http.StatusOK {
				if err != nil {
					t.Fatalf("Receive error: %s", err)
				}
				if len(got) < 1 {
					t.Fatalf("Receive %d payment methods", len(got))
				}
				return
			}

			if err != nil && errors.Is(err, tt.expectedErr) {
				t.Fatalf("Receive error: %s | must be: %s", err, tt.expectedErr)
			}

			if len(got) > 0 {
				t.Fatalf("Receive %d payment methods", len(got))
			}
		})
	}
}
