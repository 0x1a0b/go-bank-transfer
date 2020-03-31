package stub

import (
	"time"

	"github.com/gsabadini/go-bank-transfer/domain"

	"github.com/pkg/errors"
)

//TransferRepositoryStubSuccess implementa a interface de TransferRepository com resultados de sucesso
type TransferRepositoryStubSuccess struct{}

//Store cria uma transferência pré-definida
func (t TransferRepositoryStubSuccess) Store(_ domain.Transfer) (domain.Transfer, error) {
	return domain.Transfer{
		ID:                   "5e570851adcef50116aa7a5a",
		AccountOriginID:      "5e570851adcef50116aa7a5d",
		AccountDestinationID: "5e570851adcef50116aa7a5c",
		Amount:               20,
		CreatedAt:            time.Time{},
	}, nil
}

//FindAll retorna uma lista de transferências pré-definidas
func (t TransferRepositoryStubSuccess) FindAll() ([]domain.Transfer, error) {
	return []domain.Transfer{
		{
			ID:                   "5e570851adcef50116aa7a5a",
			AccountOriginID:      "5e570851adcef50116aa7a5d",
			AccountDestinationID: "5e570851adcef50116aa7a5c",
			Amount:               100,
			CreatedAt:            time.Time{},
		},
		{
			ID:                   "5e570851adcef50116aa7a5b",
			AccountOriginID:      "5e570851adcef50116aa7a5d",
			AccountDestinationID: "5e570851adcef50116aa7a5c",
			Amount:               500,
			CreatedAt:            time.Time{},
		},
	}, nil
}

//TransferRepositoryMockError implementa a interface de TransferRepository com resultados de erro
type TransferRepositoryStubError struct{}

//Store retorna um error ao criar uma transferência
func (t TransferRepositoryStubError) Store(_ domain.Transfer) (domain.Transfer, error) {
	return domain.Transfer{}, errors.New("Error")
}

//FindAll retorna um error ao listar as transferências
func (t TransferRepositoryStubError) FindAll() ([]domain.Transfer, error) {
	return []domain.Transfer{}, errors.New("Error")
}
