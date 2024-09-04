// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.16.0

package models

import (
	"context"
)

type Querier interface {
	AccountExists(ctx context.Context, uuid string) (bool, error)
	CreateAccount(ctx context.Context, arg CreateAccountParams) (*Account, error)
	CreateTransaction(ctx context.Context, arg CreateTransactionParams) (*Transaction, error)
	GetAccountDetailsByUUID(ctx context.Context, uuid string) (*Account, error)
	GetOperationTypeAmountBehavior(ctx context.Context, serialID int64) (AmountBehavior, error)
	GetTransactionDetailsByTransactionId(ctx context.Context, uuid string) (*Transaction, error)
	UserExists(ctx context.Context, uuid string) (bool, error)
}

var _ Querier = (*Queries)(nil)
