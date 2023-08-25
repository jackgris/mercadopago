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

func TestGetTestUser(t *testing.T) {

	var rightResponse string = `
         {
	   "id": 1468493374,
	   "nickname": "TESTUSER1122246296",
	   "password": "PSRt2aBZ89",
	   "site_status": "active",
	   "site_id": "",
	   "description": "",
	   "email": "test_user_7893726298@testuser.com",
	   "date_created": "",
	   "date_last_updated": ""
         }
        `
	var response mercadopago.TestUser
	if err := json.NewDecoder(strings.NewReader(rightResponse)).Decode(&response); err != nil {
		t.Fatal(err)
	}
	accessToken := "TEST-7237123416497470-080318-abc3babd65d6d886dd1193889f2b85a4-470823344"
	ctx := context.Background()

	tests := []struct {
		name             string
		respStatus       int
		accessToken      string
		siteID           string
		description      string
		expectedResponse *mercadopago.TestUser
		expectedErr      *mercadopago.ErrorResponse
	}{
		{
			name:             "Successful response",
			respStatus:       http.StatusOK,
			accessToken:      accessToken,
			siteID:           "MLA",
			description:      "This is a new user",
			expectedResponse: &response,
			expectedErr:      nil,
		},
		{
			name:             "Invalid Client ID",
			respStatus:       http.StatusBadRequest,
			accessToken:      accessToken + "longer",
			siteID:           "MLA",
			description:      "This is a new user",
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
			siteID:           "MLA",
			description:      "This is a new user",
			expectedResponse: nil,
			expectedErr: &mercadopago.ErrorResponse{
				Message: "invalid access token",
				Errors:  "",
				Status:  http.StatusBadRequest,
			},
		},
		{
			name:             "Invalid Site ID",
			respStatus:       http.StatusBadRequest,
			accessToken:      accessToken,
			siteID:           "SLA",
			description:      "This is a new user",
			expectedResponse: &response,
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

				urlPAth := "/users/test_user"
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

				var reqTestUser mercadopago.ResquestTestUser
				if err := json.NewDecoder(r.Body).Decode(&reqTestUser); err != nil {
					_, _ = w.Write([]byte("Can't parse response"))
					return
				}

				if reqTestUser.SiteID != tt.siteID {
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
			got, err := client.GetTestUser(ctx, tt.accessToken, tt.siteID, tt.description)
			if err != nil && errors.Is(err, tt.expectedErr) {
				t.Fatalf("Receive error: %s | must be: %s", err, tt.expectedErr)
			}
			if !reflect.DeepEqual(got, tt.expectedResponse) {
				t.Fatalf("Expected result is %v but receive %v", tt.expectedResponse, got)
			}
		})
	}

}
