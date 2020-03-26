package middleware

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"regexp"

	"github.com/gsabadini/go-bank-transfer/api/action"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

//ValidateAccount armazena a estrutura de validação de entrada de dados
type ValidateAccount struct {
	logger *logrus.Logger
}

//NewValidateAccount constrói um ValidateAccount com suas dependências
func NewValidateAccount(log *logrus.Logger) ValidateAccount {
	return ValidateAccount{logger: log}
}

//Execute válida os dados de criação de conta
func (v ValidateAccount) Execute(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	const (
		logKey              = "validate_account_middleware"
		messageInvalidField = "Invalid field"
	)

	payload, err := ioutil.ReadAll(r.Body)
	if err != nil {
		v.logger.WithFields(logrus.Fields{
			"key":         logKey,
			"http_status": http.StatusBadRequest,
			"error":       err.Error(),
		}).Error("error read body")

		return
	}

	// re-adiciona o payload ao buffer da request
	r.Body = ioutil.NopCloser(bytes.NewBuffer(payload))

	var account accountRequest
	if err := json.NewDecoder(r.Body).Decode(&account); err != nil {
		v.logger.WithFields(logrus.Fields{
			"key":         logKey,
			"http_status": http.StatusBadRequest,
			"error":       err.Error(),
		}).Error("error when decoding json")

		action.ErrorMessage(err, http.StatusBadRequest).Send(w)
		return
	}

	if err := account.validateBalance(); err != nil {
		v.logger.WithFields(logrus.Fields{
			"key":         logKey,
			"http_status": http.StatusBadRequest,
			"error":       err.Error(),
		}).Error(messageInvalidField)

		action.ErrorMessage(err, http.StatusBadRequest).Send(w)
		return
	}

	if err := account.validateCPF(); err != nil {
		v.logger.WithFields(logrus.Fields{
			"key":         logKey,
			"http_status": http.StatusBadRequest,
			"error":       err.Error(),
		}).Error(messageInvalidField)

		action.ErrorMessage(err, http.StatusBadRequest).Send(w)
		return
	}

	if err := account.validateName(); err != nil {
		v.logger.WithFields(logrus.Fields{
			"key":         logKey,
			"http_status": http.StatusBadRequest,
			"error":       err.Error(),
		}).Error(messageInvalidField)

		action.ErrorMessage(err, http.StatusBadRequest).Send(w)
		return
	}

	// re-adiciona o payload ao buffer da request
	r.Body = ioutil.NopCloser(bytes.NewBuffer(payload))

	next.ServeHTTP(w, r)
}

var (
	errBalanceInvalid = errors.New("Balance invalid")
	errCPFInvalid     = errors.New("CPF invalid")
	errNameInvalid    = errors.New("Name invalid")
)

type accountRequest struct {
	Name    string  `json:"name"`
	CPF     string  `json:"cpf"`
	Balance float64 `json:"balance"`
}

func (a *accountRequest) validateBalance() error {
	if a.Balance < 0 {
		return errBalanceInvalid
	}

	return nil
}

func (a *accountRequest) validateCPF() error {
	var CPFRegexp = regexp.MustCompile(`^\d{3}\.?\d{3}\.?\d{3}-?\d{2}$`)

	if !CPFRegexp.MatchString(a.CPF) {
		return errCPFInvalid
	}

	return nil
}

func (a *accountRequest) validateName() error {
	if a.Name == "" {
		return errNameInvalid
	}

	return nil
}
