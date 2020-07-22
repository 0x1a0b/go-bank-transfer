package usecase

import (
	"time"

	"github.com/gsabadini/go-bank-transfer/domain"
)

type CreateAccountInteractor struct {
	repo domain.AccountRepository
}

func NewCreateAccountInteractor(repo domain.AccountRepository) *CreateAccountInteractor {
	return &CreateAccountInteractor{repo: repo}
}

//Store cria uma nova Account
func (a CreateAccountInteractor) Store(name, CPF string, balance float64) (AccountOutput, error) {
	var account = domain.NewAccount(
		domain.AccountID(domain.NewUUID()),
		name,
		CPF,
		balance,
		time.Now(),
	)

	account, err := a.repo.Store(account)
	if err != nil {
		return AccountOutput{}, err
	}

	return AccountOutput{
		ID:        string(account.ID),
		Name:      account.Name,
		CPF:       account.CPF,
		Balance:   account.Balance,
		CreatedAt: account.CreatedAt,
	}, nil
}

