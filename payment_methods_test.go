package mercadopago_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"testing"

	"github.com/jackgris/mercadopago"
)

func TestGetPaymentMethods(t *testing.T) {

	// Get mocking data
	jsonData, err := os.Open("./testdata/payment_methods_data.json")
	if err != nil {
		t.Fatal(err)
	}
	defer jsonData.Close()
	paymentMethodsMock, _ := io.ReadAll(jsonData)

	var paymentMethods mercadopago.PaymentMethods
	if err := json.NewDecoder(bytes.NewReader(paymentMethodsMock)).Decode(&paymentMethods); err != nil {
		t.Fatal(err)
	}

	accessToken := "TEST-7237123416497470-080318-abc3babd65d6d886dd1193889f2b85a4-470823344"
	ctx := context.Background()
	data, _ := json.Marshal(paymentMethods)

	tests := []struct {
		name             string
		respStatus       int
		accessToken      string
		expectedResponse mercadopago.PaymentMethods
		expectedErr      *mercadopago.ErrorResponse
		respBody         []byte
	}{
		{
			name:             "Successful response",
			respStatus:       http.StatusOK,
			accessToken:      accessToken,
			expectedResponse: paymentMethods,
			expectedErr:      nil,
			respBody:         data,
		},
		{
			name:             "Unauthorized",
			respStatus:       http.StatusUnauthorized,
			accessToken:      "",
			expectedResponse: nil,
			expectedErr: &mercadopago.ErrorResponse{
				Message: "neither a public key or caller id were provided",
				Errors:  "unauthorized_scopes",
				Status:  http.StatusUnauthorized,
			},
			respBody: data,
		},
		{
			name:             "Not Found",
			respStatus:       http.StatusUnauthorized,
			accessToken:      "TEST-1237334875597470-080318-abc3babd65d6d886dd1193889f2b85a4-4708",
			expectedResponse: nil,
			expectedErr: &mercadopago.ErrorResponse{
				Message: "invalid_token",
				Errors:  "not_found",
				Status:  http.StatusUnauthorized,
			},
			respBody: data,
		},
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {

			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if http.MethodGet != r.Method {
					t.Fatalf("HTTP Method expected is %s but receive %s", http.MethodGet, r.Method)
				}
				urlPAth := "/v1/payment_methods"
				if urlPAth != r.URL.Path {
					t.Fatalf("Path URL expected is %s but receive %s", urlPAth, r.URL.Path)
				}

				w.Header().Set("Content-Type", "application/json; charset=utf-8")
				expectedToken := "Bearer " + tt.accessToken
				receivedToekn := r.Header.Get("Authorization")
				if len(receivedToekn) < 77 {
					w.WriteHeader(tt.respStatus)
					data, _ := json.Marshal(tt.expectedErr)
					_, _ = w.Write(data)
					return
				}

				if expectedToken != receivedToekn {
					t.Fatalf("Authorization token expected is %s but receive %s", expectedToken, receivedToekn)
				}

				w.WriteHeader(http.StatusOK)
				_, _ = w.Write(tt.respBody)
			}))

			client := mercadopago.NewClient(server.URL+"/", tt.accessToken)
			got, err := client.PaymentMethods(ctx)
			if err != nil && errors.Is(err, tt.expectedErr) {
				t.Fatalf("Receive error: %s | must be: %s", err, tt.expectedErr)
			}

			if !reflect.DeepEqual(got, tt.expectedResponse) {
				t.Fatalf("Expected result is %v but receive %v", tt.expectedResponse, got)
			}
		})
	}
}
