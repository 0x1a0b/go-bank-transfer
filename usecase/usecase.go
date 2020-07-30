package usecase

import (
	"context"

	"github.com/gsabadini/go-bank-transfer/domain"
)

//AccountUseCase é uma abstração para os casos de uso de Account
type AccountUseCase interface {
	Store(context.Context, string, string, float64) (AccountOutput, error)
	FindAll(context.Context) ([]AccountOutput, error)
	FindBalance(context.Context, domain.AccountID) (AccountBalanceOutput, error)
}

//TransferUseCase é uma abstração para os casos de uso de Transfer
type TransferUseCase interface {
	Store(context.Context, domain.AccountID, domain.AccountID, float64) (TransferOutput, error)
	FindAll(context.Context) ([]TransferOutput, error)
}
