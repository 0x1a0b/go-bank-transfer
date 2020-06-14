package middleware

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gsabadini/go-bank-transfer/mock"

	"github.com/gorilla/mux"
)

func TestValidateTransfer_Execute(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name               string
		rawPayload         []byte
		expectedStatusCode int
	}{
		{
			name: "Valid transfer",
			rawPayload: []byte(
				`{
					"account_destination_id": "3c096a40-ccba-4b58-93ed-57379ab04680" ,
					"account_origin_id": "3c096a40-ccba-4b58-93ed-57379ab04681",
					"amount": 1.00
				}`,
			),
			expectedStatusCode: http.StatusOK,
		},
		{
			name: "Invalid JSON",
			rawPayload: []byte(
				`{
					"account_destination_id": ,
					"account_origin_id": ,
					"amount": 1.00
				}`,
			),
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name: "Invalid amount",
			rawPayload: []byte(
				`{
					"account_destination_id": "3c096a40-ccba-4b58-93ed-57379ab04680",
					"account_origin_id": "3c096a40-ccba-4b58-93ed-57379ab04681",
					"amount": -1.00
				}`,
			),
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name: "Invalid account origin equals destination",
			rawPayload: []byte(
				`{
					"account_destination_id": "3c096a40-ccba-4b58-93ed-57379ab04680",
					"account_origin_id": "3c096a40-ccba-4b58-93ed-57379ab04680",
					"amount": 1.00
				}`,
			),
			expectedStatusCode: http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var body = bytes.NewReader(tt.rawPayload)
			req, err := http.NewRequest(http.MethodPost, "/transfers", body)
			if err != nil {
				t.Fatal(err)
			}

			// transformando middleware em um http.Handler
			middlewareHandler := func(w http.ResponseWriter, r *http.Request) {
				next := func(w http.ResponseWriter, r *http.Request) {}
				middleware := NewValidateTransfer(mock.LoggerMock{})
				middleware.Execute(w, r, next)
			}

			var (
				rr = httptest.NewRecorder()
				r  = mux.NewRouter()
			)

			r.HandleFunc("/transfers", middlewareHandler).Methods(http.MethodPost)
			r.ServeHTTP(rr, req)

			if rr.Code != tt.expectedStatusCode {
				t.Errorf(
					"[TestCase '%s'] O handler retornou um HTTP status code inesperado: retornado '%v' esperado '%v'",
					tt.name,
					rr.Code,
					tt.expectedStatusCode,
				)
			}
		})
	}
}
