package mercadopago_test

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/jackgris/mercadopago"
)

func TestGetCreditCard(t *testing.T) {

	var rightResponse string = `
         {
            "id": "234bc93745ac08281fdc7dba4aa4567b",
            "first_six_digits": "503123",
            "expiration_month": 10,
            "expiration_year": 2024,
            "last_four_digits": "0704",
            "cardholder": {
              "identification": {}
            },
            "status": "active",
            "date_created": "2023-08-29T23:11:19.758-04:00",
            "date_last_updated": "2023-08-29T23:11:19.758-04:00",
            "date_due": "2023-09-06T23:11:19.758-04:00",
            "luhn_validation": true,
            "live_mode": true,
            "require_esc": false,
            "card_number_length": 16,
            "security_code_length": 3
         }
        `
	var response mercadopago.CardToken
	if err := json.NewDecoder(strings.NewReader(rightResponse)).Decode(&response); err != nil {
		t.Fatal(err)
	}
	validCardToken := mercadopago.RequestCardToken{
		CardNumber:      "5031235734530705",
		ExpirationMonth: 10,
		ExpirationYear:  2024,
		SecurityCode:    "123",
		Cardholder: struct {
			Name string `json:"name"`
		}{Name: ""},
	}

	invalidCardToken := mercadopago.RequestCardToken{
		CardNumber:      "503123573453",
		ExpirationMonth: 18,
		ExpirationYear:  2000,
		SecurityCode:    "1234",
		Cardholder: struct {
			Name string `json:"name"`
		}{Name: ""},
	}

	accessToken := "TEST-7237123416497470-080318-abc3babd65d6d886dd1193889f2b85a4-470823344"
	ctx := context.Background()

	tests := []struct {
		name             string
		respStatus       int
		accessToken      string
		requestCard      *mercadopago.RequestCardToken
		expectedResponse *mercadopago.CardToken
		expectedErr      *mercadopago.ErrorResponse
	}{
		{
			name:             "Successful response",
			respStatus:       http.StatusOK,
			accessToken:      accessToken,
			requestCard:      &validCardToken,
			expectedResponse: &response,
			expectedErr:      nil,
		},
		{
			name:             "Invalid Client ID",
			respStatus:       http.StatusBadRequest,
			accessToken:      accessToken + "longer",
			requestCard:      &mercadopago.RequestCardToken{},
			expectedResponse: nil,
			expectedErr: &mercadopago.ErrorResponse{
				Message: "Invalid Client Id",
				Errors:  "bad_request",
				Status:  http.StatusBadRequest,
			},
		},
		{
			name:             "Invalid Access Token",
			respStatus:       http.StatusBadRequest,
			accessToken:      "TEST-7237123416497470-080318-abc3babd65d6d886dd1193889f2b85a4-4708233",
			requestCard:      &mercadopago.RequestCardToken{},
			expectedResponse: nil,
			expectedErr: &mercadopago.ErrorResponse{
				Message: "invalid access token",
				Errors:  "",
				Status:  http.StatusBadRequest,
			},
		},
		{
			name:             "Invalid Credit Card",
			respStatus:       http.StatusBadRequest,
			accessToken:      accessToken,
			requestCard:      &invalidCardToken,
			expectedResponse: nil,
			expectedErr: &mercadopago.ErrorResponse{
				Message: "Must be a valid site_id, please check: SLA",
				Errors:  "invalid_data",
				Status:  http.StatusBadRequest,
			},
		},
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {

			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

				if http.MethodPost != r.Method {
					errorRes := mercadopago.ErrorResponse{
						Message: fmt.Sprintf("HTTP Method expected is %s but receive %s", http.MethodPost, r.Method),
						Errors:  "Method Not Allowed",
						Status:  http.StatusMethodNotAllowed,
					}
					w.WriteHeader(http.StatusMethodNotAllowed)
					data, _ := json.Marshal(errorRes)
					_, _ = w.Write(data)
					return
				}

				urlPAth := "/v1/card_tokens"
				if urlPAth != r.URL.Path {
					errorRes := mercadopago.ErrorResponse{
						Message: fmt.Sprintf("Path URL expected is %s but receive %s", urlPAth, r.URL.Path),
						Errors:  "Wrong URL Path",
						Status:  http.StatusBadRequest,
					}
					w.WriteHeader(http.StatusBadRequest)
					data, _ := json.Marshal(errorRes)
					_, _ = w.Write(data)
					return
				}

				w.Header().Set("Content-Type", "application/json; charset=utf-8")
				expectedToken := "Bearer " + accessToken
				receivedToken := r.Header.Get("Authorization")

				if expectedToken != receivedToken {
					w.WriteHeader(tt.respStatus)
					data, _ := json.Marshal(tt.expectedErr)
					_, _ = w.Write(data)
					return
				}

				var reqCard mercadopago.RequestCardToken
				if err := json.NewDecoder(r.Body).Decode(&reqCard); err != nil {
					_, _ = w.Write([]byte("Can't parse response"))
					return
				}

				if len(reqCard.CardNumber) != response.CardNumberLength {
					w.WriteHeader(tt.respStatus)
					data, _ := json.Marshal(tt.expectedErr)
					_, _ = w.Write(data)
					return
				}

				w.WriteHeader(http.StatusOK)
				data, _ := json.Marshal(response)
				_, _ = w.Write(data)
			}))

			client := mercadopago.NewClient(server.URL+"/", tt.accessToken)
			got, err := client.GetCardToken(ctx, *tt.requestCard)
			if err != nil && errors.Is(err, tt.expectedErr) {
				t.Fatalf("Receive error: %s | must be: %s", err, tt.expectedErr)
			}
			if !reflect.DeepEqual(got, tt.expectedResponse) {
				t.Fatalf("Expected result is %v but receive %v", tt.expectedResponse, got)
			}
		})
	}

}
