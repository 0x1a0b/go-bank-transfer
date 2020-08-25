package usecase

import (
	"context"

	"github.com/gsabadini/go-bank-transfer/domain"
)

//AccountUseCase é uma abstração para os casos de uso de Account
type AccountUseCase interface {
	Store(context.Context, string, string, domain.Money) (AccountOutput, error)
	FindAll(context.Context) ([]AccountOutput, error)
	FindBalance(context.Context, domain.AccountID) (AccountBalanceOutput, error)
}
