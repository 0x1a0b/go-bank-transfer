package repository

import (
	"github.com/gsabadini/go-bank-transfer/domain"
	"github.com/gsabadini/go-bank-transfer/infrastructure/database"

	"github.com/pkg/errors"
)

const transfersCollectionName = "transfers"

//Transfer representa um repositório para dados de transferência
type Transfer DbRepository

//NewTransfer cria um repository com suas dependências
func NewTransfer(dbHandler database.NoSQLDBHandler) Transfer {
	return Transfer{dbHandler: dbHandler}
}

//Store cria uma transferência através da implementação real do database
func (t Transfer) Store(transfer domain.Transfer) (domain.Transfer, error) {
	if err := t.dbHandler.Store(transfersCollectionName, &transfer); err != nil {
		return domain.Transfer{}, errors.Wrap(err, "error creating transfer")
	}

	return transfer, nil
}

//FindAll realiza uma busca através da implementação real do database
func (t Transfer) FindAll() ([]domain.Transfer, error) {
	var transfer = make([]domain.Transfer, 0)

	if err := t.dbHandler.FindAll(transfersCollectionName, nil, &transfer); err != nil {
		return []domain.Transfer{}, errors.Wrap(err, "error listing transfers")
	}

	return transfer, nil
}
