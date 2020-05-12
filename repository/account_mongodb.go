package repository

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/gsabadini/go-bank-transfer/domain"
	"github.com/gsabadini/go-bank-transfer/infrastructure/database"

	"github.com/pkg/errors"
)

//AccountMongoDB representa um repositório para manipulação de dados de contas utilizando MongoDB
type AccountMongoDB struct {
	handler        database.NoSQLHandler
	collectionName string
}

//NewAccountMongoDB constrói um repository com suas dependências
func NewAccountMongoDB(dbHandler database.NoSQLHandler) AccountMongoDB {
	return AccountMongoDB{handler: dbHandler, collectionName: "accounts"}
}

//Store realiza a inserção de uma conta no banco de dados
func (a AccountMongoDB) Store(account domain.Account) (domain.Account, error) {
	if err := a.handler.Store(a.collectionName, &account); err != nil {
		return domain.Account{}, errors.Wrap(err, "error creating account")
	}

	return account, nil
}

//UpdateBalance realiza a atualização do saldo de uma conta no banco de dados
func (a AccountMongoDB) UpdateBalance(ID string, balance float64) error {
	var (
		query  = bson.M{"id": ID}
		update = bson.M{"$set": bson.M{"balance": balance}}
	)

	if err := a.handler.Update(a.collectionName, query, update); err != nil {
		return errors.Wrap(domain.ErrNotFound, "error updating account balance")
	}

	return nil
}

//FindAll realiza a busca de todas as contas no banco de dados
func (a AccountMongoDB) FindAll() ([]domain.Account, error) {
	var accounts = make([]domain.Account, 0)

	if err := a.handler.FindAll(a.collectionName, nil, &accounts); err != nil {
		return accounts, errors.Wrap(err, "error listing accounts")
	}

	return accounts, nil
}

//FindByID realiza a busca de uma conta no banco de dados
func (a AccountMongoDB) FindByID(ID string) (*domain.Account, error) {
	var (
		account = &domain.Account{}
		query   = bson.M{"id": ID}
	)

	if err := a.handler.FindOne(a.collectionName, query, nil, &account); err != nil {
		switch err {
		case mgo.ErrNotFound:
			return account, errors.Wrap(domain.ErrNotFound, "error fetching account")
		default:
			return account, errors.Wrap(err, "error fetching account")
		}
	}

	return account, nil
}

//FindBalance realiza a busca do saldo de uma conta no banco de dados
func (a AccountMongoDB) FindBalance(ID string) (domain.Account, error) {
	var (
		account  = domain.Account{}
		query    = bson.M{"id": ID}
		selector = bson.M{"balance": 1, "_id": 0}
	)

	if err := a.handler.FindOne(a.collectionName, query, selector, &account); err != nil {
		switch err {
		case mgo.ErrNotFound:
			return account, errors.Wrap(domain.ErrNotFound, "error fetching account balance")
		default:
			return account, errors.Wrap(err, "error fetching account balance")
		}
	}

	return account, nil
}
