package action

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/gsabadini/go-bank-transfer/domain"
	"github.com/gsabadini/go-bank-transfer/infrastructure/logger"
	"github.com/gsabadini/go-bank-transfer/infrastructure/validator"
	"github.com/gsabadini/go-bank-transfer/usecase"
	"github.com/gsabadini/go-bank-transfer/usecase/input"
	"github.com/gsabadini/go-bank-transfer/usecase/output"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

type mockCreateTransfer struct {
	result output.TransferOutput
	err    error
}

func (m mockCreateTransfer) Execute(_ context.Context, _ input.Transfer) (output.TransferOutput, error) {
	return m.result, m.err
}

func TestCreateTransferAction_Execute(t *testing.T) {
	t.Parallel()

	validator, _ := validator.NewValidatorFactory(validator.InstanceGoPlayground)

	type args struct {
		rawPayload []byte
	}

	tests := []struct {
		name               string
		args               args
		ucMock             usecase.CreateTransfer
		expectedBody       []byte
		expectedStatusCode int
	}{
		{
			name: "CreateTransferAction success",
			args: args{
				rawPayload: []byte(`{
					"account_destination_id": "3c096a40-ccba-4b58-93ed-57379ab04680",
					"account_origin_id": "3c096a40-ccba-4b58-93ed-57379ab04681",
					"amount": 10
				}`),
			},
			ucMock: mockCreateTransfer{
				result: output.TransferOutput{
					ID:                   "3c096a40-ccba-4b58-93ed-57379ab04679",
					AccountOriginID:      "3c096a40-ccba-4b58-93ed-57379ab04680",
					AccountDestinationID: "3c096a40-ccba-4b58-93ed-57379ab04681",
					Amount:               10,
					CreatedAt:            time.Time{},
				},
				err: nil,
			},
			expectedBody:       []byte(`{"id":"3c096a40-ccba-4b58-93ed-57379ab04679","account_origin_id":"3c096a40-ccba-4b58-93ed-57379ab04680","account_destination_id":"3c096a40-ccba-4b58-93ed-57379ab04681","amount":10,"created_at":"0001-01-01T00:00:00Z"}`),
			expectedStatusCode: http.StatusCreated,
		},
		{
			name: "CreateTransferAction generic error",
			args: args{
				rawPayload: []byte(
					`{
						"account_destination_id": "3c096a40-ccba-4b58-93ed-57379ab04680",
						"account_origin_id": "3c096a40-ccba-4b58-93ed-57379ab04681",
						"amount": 10
					}`,
				),
			},
			ucMock: mockCreateTransfer{
				result: output.TransferOutput{},
				err:    errors.New("error"),
			},
			expectedBody:       []byte(`{"errors":["error"]}`),
			expectedStatusCode: http.StatusInternalServerError,
		},
		{
			name: "CreateTransferAction error insufficient balance",
			args: args{
				rawPayload: []byte(
					`{
						"account_destination_id": "3c096a40-ccba-4b58-93ed-57379ab04680",
						"account_origin_id": "3c096a40-ccba-4b58-93ed-57379ab04681",
						"amount": 10
					}`,
				),
			},
			ucMock: mockCreateTransfer{
				result: output.TransferOutput{},
				err:    domain.ErrInsufficientBalance,
			},
			expectedBody:       []byte(`{"errors":["origin validator does not have sufficient balance"]}`),
			expectedStatusCode: http.StatusUnprocessableEntity,
		},
		//{
		//	name: "CreateTransferAction error validator origin equals validator destination",
		//	args: args{
		//		rawPayload: []byte(
		//			`{
		//				"account_destination_id": "3c096a40-ccba-4b58-93ed-57379ab04680",
		//				"account_origin_id": "3c096a40-ccba-4b58-93ed-57379ab04680",
		//				"amount": 10
		//			}`,
		//		),
		//	},
		//	ucMock: mockCreateTransfer{
		//		result: output.TransferOutput{},
		//		err:    nil,
		//	},
		//	expectedBody:       []byte(`{"errors":["validator origin equals destination validator"]}`),
		//	expectedStatusCode: http.StatusBadRequest,
		//},
		{
			name: "CreateTransferAction error invalid JSON",
			args: args{
				rawPayload: []byte(
					`{
						"account_destination_id": ,
						"account_origin_id": ,
						"amount":
					}`,
				),
			},
			ucMock: mockCreateTransfer{
				result: output.TransferOutput{},
				err:    nil,
			},
			expectedBody:       []byte(`{"errors":["invalid character ',' looking for beginning of value"]}`),
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name: "CreateTransferAction error invalid amount",
			args: args{
				rawPayload: []byte(
					`{
						"account_destination_id": "3c096a40-ccba-4b58-93ed-57379ab04680",
						"account_origin_id": "3c096a40-ccba-4b58-93ed-57379ab04681",
						"amount": -1
					}`,
				),
			},
			ucMock: mockCreateTransfer{
				result: output.TransferOutput{},
				err:    nil,
			},
			expectedBody:       []byte(`{"errors":["Amount must be greater than 0"]}`),
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name: "CreateTransferAction error invalid fields",
			args: args{
				rawPayload: []byte(
					`{
						"account_destination_id123": "3c096a40-ccba-4b58-93ed-57379ab04680",
						"account_origin_id123": "3c096a40-ccba-4b58-93ed-57379ab04681",
						"amount123": 10
					}`,
				),
			},
			ucMock: mockCreateTransfer{
				result: output.TransferOutput{},
				err:    nil,
			},
			expectedBody:       []byte(`{"errors":["AccountOriginID is a required field","AccountDestinationID is a required field","Amount must be greater than 0"]}`),
			expectedStatusCode: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest(
				http.MethodPost,
				"/transfers",
				bytes.NewReader(tt.args.rawPayload),
			)

			var (
				w      = httptest.NewRecorder()
				action = NewCreateTransferAction(tt.ucMock, logger.LoggerMock{}, validator)
			)

			action.Execute(w, req)

			if w.Code != tt.expectedStatusCode {
				t.Errorf(
					"[TestCase '%s'] O handler retornou um HTTP status code inesperado: retornado '%v' esperado '%v'",
					tt.name,
					w.Code,
					tt.expectedStatusCode,
				)
			}

			fmt.Println(w.Body.String())
			fmt.Println(string(tt.expectedBody))
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
