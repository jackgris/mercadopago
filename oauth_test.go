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

func TestGetAccessToken(t *testing.T) {

	var rightResponse string = `
        {
         "access_token": "APP_USR-7237385876897478-882419-cf8589ace9fee57cb876a2dc72ed88a6-57883988",
         "token_type": "Bearer",
         "expires_in": 21600,
         "scope": "offline_access payments read write",
         "user_id": 87883988
        }
        `
	var response mercadopago.AccessToken
	if err := json.NewDecoder(strings.NewReader(rightResponse)).Decode(&response); err != nil {
		t.Fatal(err)
	}
	accessToken := "TEST-7237123416497470-080318-abc3babd65d6d886dd1193889f2b85a4-470823344"
	ctx := context.Background()

	tests := []struct {
		name             string
		respStatus       int
		accessToken      string
		clientID         string
		clientSecret     string
		expectedResponse *mercadopago.AccessToken
		expectedErr      *mercadopago.ErrorResponse
	}{
		{
			name:             "Successful response",
			respStatus:       http.StatusOK,
			accessToken:      accessToken,
			clientID:         "9837385876897878",
			clientSecret:     "9h8WjMhqOkpaxofv8yjdMtajkoyJMm8R",
			expectedResponse: &response,
			expectedErr:      nil,
		},
		{
			name:             "Empty client id or client secret",
			respStatus:       http.StatusBadRequest,
			accessToken:      accessToken,
			clientID:         "",
			clientSecret:     "",
			expectedResponse: nil,
			expectedErr: &mercadopago.ErrorResponse{
				Message: "the following parameters are required: grant_type, client_id, client_secret. Missing parameters: client_id",
				Errors:  "invalid_request",
				Status:  http.StatusBadRequest,
			},
		},
		{
			name:             "Wront clientId or clientSecret",
			respStatus:       http.StatusBadRequest,
			accessToken:      accessToken,
			clientID:         "98373858768",
			clientSecret:     "9h8WjMhqOkpaxofv8yjdMtajkoyJ",
			expectedResponse: nil,
			expectedErr: &mercadopago.ErrorResponse{
				Message: "",
				Errors:  "invalid_client",
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

				urlPAth := "/oauth/token"
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
					errorRes := mercadopago.ErrorResponse{
						Message: fmt.Sprintf("Authorization token expected is %s but receive %s", expectedToken, receivedToekn),
						Errors:  "Not Authorized",
						Status:  http.StatusBadGateway,
					}
					w.WriteHeader(http.StatusBadRequest)
					data, _ := json.Marshal(errorRes)
					_, _ = w.Write(data)
					return
				}

				var reqAccessToken mercadopago.RequestAccessToken
				if err := json.NewDecoder(r.Body).Decode(&reqAccessToken); err != nil {
					_, _ = w.Write([]byte("Can't parse response"))
					return
				}

				if reqAccessToken.ClientID == "" || reqAccessToken.ClientSecret == "" {
					w.WriteHeader(tt.respStatus)
					data, _ := json.Marshal(tt.expectedErr)
					_, _ = w.Write(data)
					return
				}

				if len(reqAccessToken.ClientID) < 12 || len(reqAccessToken.ClientSecret) < 29 {
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
			got, err := client.GetAccessToken(ctx, tt.clientID, tt.clientSecret)
			if err != nil && errors.Is(err, tt.expectedErr) {
				t.Fatalf("Receive error: %s | must be: %s", err, tt.expectedErr)
			}
			if !reflect.DeepEqual(got, tt.expectedResponse) {
				t.Fatalf("Expected result is %v but receive %v", tt.expectedResponse, got)
			}
		})
	}

}
