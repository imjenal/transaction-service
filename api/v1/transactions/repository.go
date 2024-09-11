package transactions

import (
	"context"
	"errors"
	"fmt"

	"github.com/imjenal/transaction-service/internal/db/models"
	"github.com/jackc/pgx/v4"
)

type Repository struct {
	querier models.Querier
}

func NewRepository(querier models.Querier) *Repository {
	return &Repository{querier: querier}
}

var (
	errTransactionNotFound   = errors.New("TRANSACTION_NOT_FOUND")
	errOperationTypeNotFound = errors.New("OPERATION_TYPE_NOT_FOUND")
	errAccountNotFound       = errors.New("ACCOUNT_NOT_FOUND")
)

func (r *Repository) getTransactionDetails(ctx context.Context, uuid string) (*models.GetTransactionDetailsByTransactionIdRow, error) {
	transactionDetails, err := r.querier.GetTransactionDetailsByTransactionId(ctx, uuid)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, errTransactionNotFound
	}

	if err != nil {
		return nil, fmt.Errorf("repo.getTransactionDetails: error : %w", err)
	}

	return transactionDetails, nil
}

func (r *Repository) createTransaction(ctx context.Context, arg models.CreateTransactionParams) (*models.CreateTransactionRow, error) {
	transactionDetails, err := r.querier.CreateTransaction(ctx, arg)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}

	if err != nil {
		return nil, fmt.Errorf("repo.createTransaction: error: %w", err)
	}

	return transactionDetails, nil
}

func (r *Repository) accountExists(ctx context.Context, accountID string) (bool, error) {
	exists, err := r.querier.AccountExists(ctx, accountID)
	if err != nil {
		return false, fmt.Errorf("repo.accountExists: error checking account existence: %w", err)
	}
	return exists, nil
}

func (r *Repository) getAmountBehavior(ctx context.Context, operationTypeID int64) (models.AmountBehavior, error) {
	amountBehavior, err := r.querier.GetOperationTypeAmountBehavior(ctx, operationTypeID)
	if errors.Is(err, pgx.ErrNoRows) {
		return "", errOperationTypeNotFound
	}

	if err != nil {
		return "", fmt.Errorf("repo.getAmountBehavior: error fetching amount behavior: %w", err)
	}
	return amountBehavior, nil
}

func (r *Repository) getNegativeBalanceTransactionsByAccountID(ctx context.Context, accountID string) ([]*models.GetNegativeBalanceTransactionsByAccountIDRow, error) {
	transactions, err := r.querier.GetNegativeBalanceTransactionsByAccountID(ctx, accountID)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}

	if err != nil {
		return nil, fmt.Errorf("repo.getNegativeBalanceTransactionsByAccountID: error fetching negative balance txns: %w", err)
	}
	return transactions, nil
}

func (r *Repository) updateTransactionBalances(ctx context.Context, transactions []*models.GetNegativeBalanceTransactionsByAccountIDRow) error {
	for _, txn := range transactions {
		args := models.UpdateTransactionBalancesParams{
			Uuid:    txn.Uuid,
			Balance: txn.Balance,
		}

		err := r.querier.UpdateTransactionBalances(ctx, args)
		if errors.Is(err, pgx.ErrNoRows) {
			return nil
		}

		if err != nil {
			return fmt.Errorf("repo.updateTransactionBalances: failed to update txn balance: %w", err)
		}
	}

	return nil
}
