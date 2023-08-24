package mercadopago

import "fmt"

// ErrorResponse represent the error message that the API return
type ErrorResponse struct {
	Message string        `json:"message"`
	Errors  string        `json:"error"`
	Status  int           `json:"status"`
	Cause   []interface{} `json:"cause"`
}

func (e *ErrorResponse) Error() string {
	return fmt.Sprintf("Status code: %d - Error: %s - Message: %s", e.Status, e.Errors, e.Message)
}
