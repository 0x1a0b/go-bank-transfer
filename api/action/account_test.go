package action

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gsabadini/go-bank-transfer/domain"
	"github.com/gsabadini/go-bank-transfer/infrastructure/logger"
	"github.com/gsabadini/go-bank-transfer/infrastructure/validator"
	"github.com/gsabadini/go-bank-transfer/usecase"
)

type mockAccountStore struct {
	usecase.AccountUseCase

	result usecase.AccountOutput
	err    error
}

func (m mockAccountStore) Store(_, _ string, _ float64) (usecase.AccountOutput, error) {
	return m.result, m.err
}

func TestAccount_Store(t *testing.T) {
	t.Parallel()

	validator, _ := validator.NewValidatorFactory(validator.InstanceGoPlayground, logger.LoggerMock{})

	type args struct {
		rawPayload []byte
	}

	tests := []struct {
		name               string
		args               args
		usecaseMock        usecase.AccountUseCase
		expectedBody       []byte
		expectedStatusCode int
	}{
		{
			name: "Store action success",
			args: args{
				rawPayload: []byte(
					`{
						"name": "test",
						"cpf": "44451598087", 
						"balance": 10 
					}`,
				),
			},
			usecaseMock: mockAccountStore{
				result: usecase.AccountOutput{
					ID:        "3c096a40-ccba-4b58-93ed-57379ab04680",
					Name:      "Test",
					CPF:       "07094564964",
					Balance:   10,
					CreatedAt: time.Time{},
				},
				err: nil,
			},
			expectedBody:       []byte(`{"id":"3c096a40-ccba-4b58-93ed-57379ab04680","name":"Test","cpf":"07094564964","balance":10,"created_at":"0001-01-01T00:00:00Z"}`),
			expectedStatusCode: http.StatusCreated,
		},
		{
			name: "Store action generic error",
			args: args{
				rawPayload: []byte(
					`{
						"name": "test",
						"cpf": "44451598087",
						"balance": 10
					}`,
				),
			},
			usecaseMock: mockAccountStore{
				result: usecase.AccountOutput{},
				err:    errors.New("error"),
			},
			expectedBody:       []byte(`{"errors":["error"]}`),
			expectedStatusCode: http.StatusInternalServerError,
		},
		{
			name: "Store action error invalid balance",
			args: args{
				rawPayload: []byte(
					`{
						"name": "test",
						"cpf": "44451598087",
						"balance": -1
					}`,
				),
			},
			usecaseMock: mockAccountStore{
				result: usecase.AccountOutput{},
				err:    nil,
			},
			expectedBody:       []byte(`{"errors":["Balance must be greater than 0"]}`),
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name: "Store action error invalid fields",
			args: args{
				rawPayload: []byte(
					`{
						"name123": "test",
						"cpf1231": "44451598087",
						"balance12312": 1
					}`,
				),
			},
			usecaseMock: mockAccountStore{
				result: usecase.AccountOutput{},
				err:    nil,
			},
			expectedBody:       []byte(`{"errors":["Name is a required field","CPF is a required field","Balance must be greater than 0"]}`),
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name: "Store action error invalid JSON",
			args: args{
				rawPayload: []byte(
					`{
						"name":
					}`,
				),
			},
			usecaseMock: mockAccountStore{
				result: usecase.AccountOutput{},
				err:    nil,
			},
			expectedBody:       []byte(`{"errors":["invalid character '}' looking for beginning of value"]}`),
			expectedStatusCode: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest(
				http.MethodPost,
				"/accounts",
				bytes.NewReader(tt.args.rawPayload),
			)

			var (
				w      = httptest.NewRecorder()
				action = NewAccount(tt.usecaseMock, logger.LoggerMock{}, validator)
			)

			action.Store(w, req)

			if w.Code != tt.expectedStatusCode {
				t.Errorf(
					"[TestCase '%s'] O handler retornou um HTTP status code inesperado: retornado '%v' esperado '%v'",
					tt.name,
					w.Code,
					tt.expectedStatusCode,
				)
			}

			var result = bytes.TrimSpace(w.Body.Bytes())
			if !bytes.Equal(result, tt.expectedBody) {
				t.Errorf(
					"[TestCase '%s'] Result: '%v' | Expected: '%v'",
					tt.name,
					result,
					tt.expectedBody,
				)
			}
		})
	}
}

type mockAccountFindAll struct {
	usecase.AccountUseCase

	result []usecase.AccountOutput
	err    error
}

func (m mockAccountFindAll) FindAll() ([]usecase.AccountOutput, error) {
	return m.result, m.err
}

func TestAccount_Index(t *testing.T) {
	t.Parallel()

	validator, _ := validator.NewValidatorFactory(validator.InstanceGoPlayground, logger.LoggerMock{})

	tests := []struct {
		name               string
		usecaseMock        usecase.AccountUseCase
		expectedBody       []byte
		expectedStatusCode int
	}{
		{
			name: "Index handler success one account",
			usecaseMock: mockAccountFindAll{
				result: []usecase.AccountOutput{
					{
						ID:        "3c096a40-ccba-4b58-93ed-57379ab04680",
						Name:      "Test",
						CPF:       "07094564964",
						Balance:   10,
						CreatedAt: time.Time{},
					},
				},
				err: nil,
			},
			expectedBody:       []byte(`[{"id":"3c096a40-ccba-4b58-93ed-57379ab04680","name":"Test","cpf":"07094564964","balance":10,"created_at":"0001-01-01T00:00:00Z"}]`),
			expectedStatusCode: http.StatusOK,
		},
		{
			name: "Index handler success empty",
			usecaseMock: mockAccountFindAll{
				result: []usecase.AccountOutput{},
				err:    nil,
			},
			expectedBody:       []byte(`[]`),
			expectedStatusCode: http.StatusOK,
		},
		{
			name: "Index handler generic error",
			usecaseMock: mockAccountFindAll{
				err: errors.New("error"),
			},
			expectedBody:       []byte(`{"errors":["error"]}`),
			expectedStatusCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest(http.MethodGet, "/accounts", nil)

			var (
				w      = httptest.NewRecorder()
				action = NewAccount(tt.usecaseMock, logger.LoggerMock{}, validator)
			)

			action.Index(w, req)

			if w.Code != tt.expectedStatusCode {
				t.Errorf(
					"[TestCase '%s'] O handler retornou um HTTP status code inesperado: retornado '%v' esperado '%v'",
					tt.name,
					w.Code,
					tt.expectedStatusCode,
				)
			}

			var result = bytes.TrimSpace(w.Body.Bytes())
			if !bytes.Equal(result, tt.expectedBody) {
				t.Errorf(
					"[TestCase '%s'] Result: '%v' | Expected: '%v'",
					tt.name,
					result,
					tt.expectedBody,
				)
			}
		})
	}
}

type mockAccountFindBalance struct {
	usecase.AccountUseCase

	result usecase.AccountBalanceOutput
	err    error
}

func (m mockAccountFindBalance) FindBalance(_ string) (usecase.AccountBalanceOutput, error) {
	return m.result, m.err
}

func TestAccount_FindBalance(t *testing.T) {
	t.Parallel()

	validator, _ := validator.NewValidatorFactory(validator.InstanceGoPlayground, logger.LoggerMock{})

	type args struct {
		accountID string
	}

	tests := []struct {
		name               string
		args               args
		usecaseMock        usecase.AccountUseCase
		expectedBody       []byte
		expectedStatusCode int
	}{
		{
			name: "FindBalance action success",
			args: args{
				accountID: "3c096a40-ccba-4b58-93ed-57379ab04680",
			},
			usecaseMock: mockAccountFindBalance{
				result: usecase.AccountBalanceOutput{
					Balance: 10,
				},
				err: nil,
			},
			expectedBody:       []byte(`{"balance":10}`),
			expectedStatusCode: http.StatusOK,
		},
		{
			name: "FindBalance action generic error",
			args: args{
				accountID: "3c096a40-ccba-4b58-93ed-57379ab04680",
			},
			usecaseMock: mockAccountFindBalance{
				result: usecase.AccountBalanceOutput{},
				err:    errors.New("error"),
			},
			expectedBody:       []byte(`{"errors":["error"]}`),
			expectedStatusCode: http.StatusInternalServerError,
		},
		{
			name: "FindBalance action error parameter invalid",
			args: args{
				accountID: "error",
			},
			usecaseMock: mockAccountFindBalance{
				result: usecase.AccountBalanceOutput{},
				err:    nil,
			},
			expectedBody:       []byte(`{"errors":["parameter invalid"]}`),
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name: "FindBalance action error fetching account",
			args: args{
				accountID: "3c096a40-ccba-4b58-93ed-57379ab04680",
			},
			usecaseMock: mockAccountFindBalance{
				result: usecase.AccountBalanceOutput{},
				err:    domain.ErrNotFound,
			},
			expectedBody:       []byte(`{"errors":["not found"]}`),
			expectedStatusCode: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uri := fmt.Sprintf("/accounts/%s/balance", tt.args.accountID)
			req, _ := http.NewRequest(http.MethodGet, uri, nil)

			q := req.URL.Query()
			q.Add("account_id", tt.args.accountID)
			req.URL.RawQuery = q.Encode()

			var (
				w      = httptest.NewRecorder()
				action = NewAccount(tt.usecaseMock, logger.LoggerMock{}, validator)
			)

			action.FindBalance(w, req)

			if w.Code != tt.expectedStatusCode {
				t.Errorf(
					"[TestCase '%s'] O handler retornou um HTTP status code inesperado: retornado '%v' esperado '%v'",
					tt.name,
					w.Code,
					tt.expectedStatusCode,
				)
			}

			var result = bytes.TrimSpace(w.Body.Bytes())
			if !bytes.Equal(result, tt.expectedBody) {
				t.Errorf(
					"[TestCase '%s'] Result: '%v' | Expected: '%v'",
					tt.name,
					result,
					tt.expectedBody,
				)
			}
		})
	}
}
