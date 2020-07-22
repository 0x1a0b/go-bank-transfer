package action

import (
	"encoding/json"
	"net/http"

	"github.com/gsabadini/go-bank-transfer/api/response"
	"github.com/gsabadini/go-bank-transfer/infrastructure/logger"
	"github.com/gsabadini/go-bank-transfer/infrastructure/validator"
	"github.com/gsabadini/go-bank-transfer/usecase"
)

//Account armazena as dependências para as ações de Account
type CreateAccount struct {
	uc        usecase.CreateAccountUseCase
	log       logger.Logger
	validator validator.Validator
}

func NewCreateAccount(uc usecase.CreateAccountUseCase, log logger.Logger, validator validator.Validator) *CreateAccount {
	return &CreateAccount{uc: uc, log: log, validator: validator}
}

//Store é um handler para criação de Account
func (a CreateAccount) Store(w http.ResponseWriter, r *http.Request) {
	//const logKey = "create_account"

	var input accountInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		//a.logError(
		//	logKey,
		//	"error when decoding json",
		//	http.StatusBadRequest,
		//	err,
		//)

		response.NewError(err, http.StatusBadRequest).Send(w)
		return
	}
	defer r.Body.Close()

	if errs := a.validateInput(input); len(errs) > 0 {
		//a.logError(
		//	logKey,
		//	"invalid input",
		//	http.StatusBadRequest,
		//	errors.New("invalid input"),
		//)

		response.NewErrorMessage(errs, http.StatusBadRequest).Send(w)
		return
	}

	output, err := a.uc.Store(input.Name, input.CPF, input.Balance)
	if err != nil {
		//a.logError(
		//	logKey,
		//	"error when creating a new account",
		//	http.StatusInternalServerError,
		//	err,
		//)

		response.NewError(err, http.StatusInternalServerError).Send(w)
		return
	}
	//a.logSuccess(logKey, "success creating account", http.StatusCreated)

	response.NewSuccess(output, http.StatusCreated).Send(w)
}

func (a CreateAccount) validateInput(input accountInput) []string {
	var messages []string

	err := a.validator.Validate(input)
	if err != nil {
		for _, msg := range a.validator.Messages() {
			messages = append(messages, msg)
		}
	}

	return messages
}

