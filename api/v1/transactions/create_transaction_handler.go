package transactions

import (
	"context"
	"errors"
	"github.com/imjenal/transaction-service/internal/db/models"
	"github.com/imjenal/transaction-service/pkg/http/response"
	"log"
	"math"
	"net/http"
)

type CreateTransactionRequestData struct {
	AccountId       string  `json:"account_id" validate:"required,uuid"`
	OperationTypeId int64   `json:"operation_type_id" validate:"required"`
	Amount          float64 `json:"amount"  validate:"required,gt=0"`
}

// createTransaction handles creating a transaction
func (h *Handler) createTransaction() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		requestBody := &CreateTransactionRequestData{}
		if ok := h.reader.ReadJSONAndValidate(w, r, requestBody); !ok {
			return
		}

		ctx := r.Context()

		if !h.validateAccount(ctx, w, requestBody.AccountId) {
			return
		}

		amountBehavior, err := h.validateAndFetchOperationType(ctx, w, requestBody.OperationTypeId)
		if err != nil {
			return
		}

		requestBody.Amount = adjustAmountBasedOnOperationTypeAmountBehavior(amountBehavior, requestBody.Amount)

		h.createAndRespondTransaction(ctx, w, requestBody)
	}
}

// validateAccount checks if the account exists
func (h *Handler) validateAccount(ctx context.Context, w http.ResponseWriter, accountID string) bool {
	accountExists, err := h.repository.accountExists(ctx, accountID)
	if err != nil {
		log.Printf("validateAccount: failed to check account existence: %v", err)
		h.writer.Internal(w, &response.APIError{
			Code:    response.DefaultErrorCode,
			Message: "Failed to validate account ID.",
		})
		return false
	}

	if !accountExists {
		log.Printf("validateAccount: account %s does not exist", accountID)
		h.writer.NotFound(w, &response.APIError{
			Code:    response.ErrAccountNotFound,
			Message: errAccountNotFound.Error(),
		})
		return false
	}
	return true
}

// validateAndFetchOperationType checks if the operation type exists and retrieves its amount behavior
func (h *Handler) validateAndFetchOperationType(ctx context.Context, w http.ResponseWriter, operationTypeID int64) (models.AmountBehavior, error) {
	amountBehavior, err := h.repository.getAmountBehavior(ctx, operationTypeID)
	if errors.Is(err, errOperationTypeNotFound) {
		log.Printf("validateAndFetchOperationType: operation type %d does not exist", operationTypeID)
		h.writer.NotFound(w, &response.APIError{
			Code:    response.ErrOperationTypeNotFound,
			Message: errOperationTypeNotFound.Error(),
		})
		return "", err
	}

	if err != nil {
		log.Printf("validateAndFetchOperationType: failed to get amount behavior: %v", err)
		h.writer.Internal(w, &response.APIError{
			Code:    response.DefaultErrorCode,
			Message: "Failed to fetch operation type information.",
		})
		return "", err
	}

	return amountBehavior, nil
}

// createAndRespondTransaction creates the transaction and responds to the client
func (h *Handler) createAndRespondTransaction(ctx context.Context, w http.ResponseWriter, requestBody *CreateTransactionRequestData) {
	txnDetails, err := h.repository.createTransaction(ctx, models.CreateTransactionParams{
		AccountID:       requestBody.AccountId,
		OperationTypeID: requestBody.OperationTypeId,
		Amount:          requestBody.Amount,
	})
	if err != nil {
		log.Printf("createTransaction: failed to create transaction: %v", err)
		h.writer.Internal(w, &response.APIError{
			Code:    response.DefaultErrorCode,
			Message: "Failed to create transaction.",
		})
		return
	}

	h.writer.Ok(w, txnDetails)
}

// Adjust the amount based on the amount behavior
func adjustAmountBasedOnOperationTypeAmountBehavior(amountBehavior models.AmountBehavior, amount float64) float64 {
	switch amountBehavior {
	case models.AmountBehaviorNEGATIVE:
		return -math.Abs(amount) // Store as negative
	case models.AmountBehaviorPOSITIVE:
		return math.Abs(amount) // Store as positive
	default:
		// In case of an unexpected value, return the absolute value by default
		log.Printf("Unknown amount behavior: %v, defaulting to positive amount.", amountBehavior)
		return math.Abs(amount)
	}
}
