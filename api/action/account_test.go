package action

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gsabadini/go-bank-transfer/infrastructure/database"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus/hooks/test"
)

//@TODO REVER TESTES
func TestAccountStore(t *testing.T) {
	type args struct {
		accountAction Account
		rawPayload    []byte
	}

	var loggerMock, _ = test.NewNullLogger()

	tests := []struct {
		name               string
		expectedStatusCode int
		args               args
	}{
		{
			name:               "Store handler database success",
			expectedStatusCode: http.StatusCreated,
			args: args{
				accountAction: NewAccount(database.MongoHandlerSuccessMock{}, loggerMock),
				rawPayload:    []byte(`{"name": "test","cpf": "44451598087", "balance": 10 }`),
			},
		},
		{
			name:               "Store handler database error",
			expectedStatusCode: http.StatusInternalServerError,
			args: args{
				accountAction: NewAccount(database.MongoHandlerErrorMock{}, loggerMock),
				rawPayload:    []byte(`{"name": "test","cpf": "44451598087", "balance": 10 }`),
			},
		},
		{
			name:               "Store handler invalid JSON",
			expectedStatusCode: http.StatusBadRequest,
			args: args{
				accountAction: NewAccount(database.MongoHandlerSuccessMock{}, loggerMock),
				rawPayload:    []byte(`{"name": }`),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var body = bytes.NewReader(tt.args.rawPayload)

			req, err := http.NewRequest(http.MethodPost, "/accounts", body)
			if err != nil {
				t.Fatal(err)
			}

			var (
				rr = httptest.NewRecorder()
				r  = mux.NewRouter()
			)

			r.HandleFunc("/accounts", tt.args.accountAction.Store).Methods(http.MethodPost)
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

func TestAccountIndex(t *testing.T) {
	type args struct {
		accountAction Account
	}

	var loggerMock, _ = test.NewNullLogger()

	tests := []struct {
		name               string
		expectedStatusCode int
		args               args
	}{
		{
			name:               "Index handler database success",
			expectedStatusCode: http.StatusOK,
			args: args{
				accountAction: NewAccount(database.MongoHandlerSuccessMock{}, loggerMock),
			},
		},
		{
			name:               "Index handler database error",
			expectedStatusCode: http.StatusInternalServerError,
			args: args{
				accountAction: NewAccount(database.MongoHandlerErrorMock{}, loggerMock),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest(http.MethodGet, "/accounts", nil)
			if err != nil {
				t.Fatal(err)
			}

			var (
				rr = httptest.NewRecorder()
				r  = mux.NewRouter()
			)

			r.HandleFunc("/accounts", tt.args.accountAction.Index).Methods(http.MethodGet)
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

func TestAccountShowBalance(t *testing.T) {
	type args struct {
		accountAction Account
	}

	var loggerMock, _ = test.NewNullLogger()

	tests := []struct {
		name               string
		expectedStatusCode int
		args               args
	}{
		{
			name:               "ShowBalance handler database success",
			expectedStatusCode: http.StatusOK,
			args: args{
				accountAction: NewAccount(database.MongoHandlerSuccessMock{}, loggerMock),
			},
		},
		{
			name:               "ShowBalance handler database error",
			expectedStatusCode: http.StatusInternalServerError,
			args: args{
				accountAction: NewAccount(database.MongoHandlerErrorMock{}, loggerMock),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest(http.MethodGet, "/accounts/5e5282beba39bfc244dc4c4b/balance", nil)
			if err != nil {
				t.Fatal(err)
			}

			req = mux.SetURLVars(req, map[string]string{"account_id": "5e5282beba39bfc244dc4c4b"})

			var (
				rr = httptest.NewRecorder()
				r  = mux.NewRouter()
			)

			r.HandleFunc("/accounts/{account_id}/balance", tt.args.accountAction.ShowBalance).Methods(http.MethodGet)
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
