package accounts

import (
	"context"
	"errors"
	"fmt"
	"github.com/imjenal/transaction-service/internal/db/models"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
)

type Repository struct {
	querier models.Querier
}

func NewRepository(querier models.Querier) *Repository {
	return &Repository{querier: querier}
}

var (
	errAccountNotFound      = errors.New("ACCOUNT_NOT_FOUND")
	errAccountAlreadyExists = errors.New("ACCOUNT_ALREADY_EXISTS")
	errUserNotFound         = errors.New("USER_NOT_FOUND")
)

func (r *Repository) getAccountDetails(ctx context.Context, uuid string) (*models.Account, error) {
	accountDetails, err := r.querier.GetAccountDetailsByUUID(ctx, uuid)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, errAccountNotFound
	}

	if err != nil {
		return nil, fmt.Errorf("repo.getAccountDetails: error: %w", err)
	}

	return accountDetails, nil
}

func (r *Repository) createAccount(ctx context.Context, arg models.CreateAccountParams) (*models.Account, error) {

	accountDetails, err := r.querier.CreateAccount(ctx, arg)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == "23505" { // 23505 is a unique violation
		return nil, errAccountAlreadyExists
	}

	if err != nil {
		return nil, fmt.Errorf("repo.createAccount: error: %w", err)
	}
	return accountDetails, nil
}

func (r *Repository) userExists(ctx context.Context, userID string) (bool, error) {
	exists, err := r.querier.UserExists(ctx, userID)
	if err != nil {
		return false, fmt.Errorf("repo.createAccount: error checking user existence: %w", err)
	}
	return exists, nil
}
