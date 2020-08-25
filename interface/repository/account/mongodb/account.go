package mongodb

import (
	"context"
	"time"

	"github.com/gsabadini/go-bank-transfer/domain"
	"github.com/gsabadini/go-bank-transfer/interface/repository"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

//accountBSON armazena a estrutura de dados do MongoDB
type accountBSON struct {
	ID        string    `bson:"id"`
	Name      string    `bson:"name"`
	CPF       string    `bson:"cpf"`
	Balance   int64     `bson:"balance"`
	CreatedAt time.Time `bson:"created_at"`
}

//AccountRepository armazena a estrutura de dados de um repositório de Account
type AccountRepository struct {
	collectionName string
	handler        repository.NoSQLHandler
}

//NewAccountRepository constrói um repository com suas dependências
func NewAccountRepository(h repository.NoSQLHandler) AccountRepository {
	return AccountRepository{handler: h, collectionName: "accounts"}
}

//Store insere uma Account no database
func (a AccountRepository) Store(ctx context.Context, account domain.Account) (domain.Account, error) {
	var accountBSON = accountBSON{
		ID:        account.ID().String(),
		Name:      account.Name(),
		CPF:       account.CPF(),
		Balance:   account.Balance().Int64(),
		CreatedAt: account.CreatedAt(),
	}

	if err := a.handler.Store(ctx, a.collectionName, accountBSON); err != nil {
		return domain.Account{}, errors.Wrap(err, "error creating account")
	}

	return account, nil
}

//UpdateBalance atualiza o Balance de uma Account no database
func (a AccountRepository) UpdateBalance(ctx context.Context, ID domain.AccountID, balance domain.Money) error {
	var (
		query  = bson.M{"id": ID}
		update = bson.M{"$set": bson.M{"balance": balance}}
	)

	if err := a.handler.Update(ctx, a.collectionName, query, update); err != nil {
		switch err {
		case mongo.ErrNilDocument:
			return errors.Wrap(domain.ErrNotFound, "error updating account balance")
		default:
			return errors.Wrap(err, "error updating account balance")
		}
	}

	return nil
}

//FindAll busca todas as Account no database
func (a AccountRepository) FindAll(ctx context.Context) ([]domain.Account, error) {
	var accountsBSON = make([]accountBSON, 0)

	if err := a.handler.FindAll(ctx, a.collectionName, bson.M{}, &accountsBSON); err != nil {
		switch err {
		case mongo.ErrNilDocument:
			return []domain.Account{}, errors.Wrap(domain.ErrNotFound, "error listing accounts")
		default:
			return []domain.Account{}, errors.Wrap(err, "error listing accounts")
		}
	}

	var accounts = make([]domain.Account, 0)

	for _, accountBSON := range accountsBSON {
		var account = domain.NewAccount(
			domain.AccountID(accountBSON.ID),
			accountBSON.Name,
			accountBSON.CPF,
			domain.Money(accountBSON.Balance),
			accountBSON.CreatedAt,
		)

		accounts = append(accounts, account)
	}

	return accounts, nil
}

//FindByID busca uma Account por id no database
func (a AccountRepository) FindByID(ctx context.Context, ID domain.AccountID) (domain.Account, error) {
	var (
		accountBSON = &accountBSON{}
		query       = bson.M{"id": ID}
	)

	if err := a.handler.FindOne(ctx, a.collectionName, query, nil, accountBSON); err != nil {
		switch err {
		case mongo.ErrNoDocuments:
			return domain.Account{}, errors.Wrap(domain.ErrNotFound, "error fetching account")
		default:
			return domain.Account{}, errors.Wrap(err, "error fetching account")
		}
	}

	return domain.NewAccount(
		domain.AccountID(accountBSON.ID),
		accountBSON.Name,
		accountBSON.CPF,
		domain.Money(accountBSON.Balance),
		accountBSON.CreatedAt,
	), nil
}

//FindBalance busca o Balance de uma Account no database
func (a AccountRepository) FindBalance(ctx context.Context, ID domain.AccountID) (domain.Account, error) {
	var (
		accountBSON = &accountBSON{}
		query       = bson.M{"id": ID}
		projection  = bson.M{"balance": 1, "_id": 0}
	)

	if err := a.handler.FindOne(ctx, a.collectionName, query, projection, accountBSON); err != nil {
		switch err {
		case mongo.ErrNoDocuments:
			return domain.Account{}, errors.Wrap(domain.ErrNotFound, "error fetching account balance")
		default:
			return domain.Account{}, errors.Wrap(err, "error fetching account balance")
		}
	}

	return domain.NewAccountBalance(domain.Money(accountBSON.Balance)), nil
}
