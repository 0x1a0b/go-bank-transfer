package mongodb

import (
	"context"
	"time"

	"github.com/gsabadini/go-bank-transfer/domain"
	"github.com/gsabadini/go-bank-transfer/interface/repository"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
)

type transferBSON struct {
	ID                   string    `bson:"id"`
	AccountOriginID      string    `bson:"account_origin_id"`
	AccountDestinationID string    `bson:"account_destination_id"`
	Amount               int64     `bson:"amount"`
	CreatedAt            time.Time `bson:"created_at"`
}

type TransferRepository struct {
	collectionName string
	handler        repository.NoSQLHandler
}

func NewTransferRepository(h repository.NoSQLHandler) TransferRepository {
	return TransferRepository{handler: h, collectionName: "transfers"}
}

func (t TransferRepository) Create(ctx context.Context, transfer domain.Transfer) (domain.Transfer, error) {
	transferBSON := &transferBSON{
		ID:                   transfer.ID().String(),
		AccountOriginID:      transfer.AccountOriginID().String(),
		AccountDestinationID: transfer.AccountDestinationID().String(),
		Amount:               transfer.Amount().Int64(),
		CreatedAt:            transfer.CreatedAt(),
	}

	if err := t.handler.Store(ctx, t.collectionName, transferBSON); err != nil {
		return domain.Transfer{}, errors.Wrap(err, "error creating transfer")
	}

	return transfer, nil
}

func (t TransferRepository) FindAll(ctx context.Context) ([]domain.Transfer, error) {
	var transfersBSON = make([]transferBSON, 0)

	if err := t.handler.FindAll(ctx, t.collectionName, bson.M{}, &transfersBSON); err != nil {
		return []domain.Transfer{}, errors.Wrap(err, "error listing transfers")
	}

	var transfers = make([]domain.Transfer, 0)

	for _, transferBSON := range transfersBSON {
		var transfer = domain.NewTransfer(
			domain.TransferID(transferBSON.ID),
			domain.AccountID(transferBSON.AccountOriginID),
			domain.AccountID(transferBSON.AccountDestinationID),
			domain.Money(transferBSON.Amount),
			transferBSON.CreatedAt,
		)

		transfers = append(transfers, transfer)
	}

	return transfers, nil
}