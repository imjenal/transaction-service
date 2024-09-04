package accounts

import (
	"context"
	"github.com/imjenal/transaction-service/internal/db/models"
	"github.com/imjenal/transaction-service/pkg/http/response"
	"log"
	"net/http"
)

type CreateAccountRequestData struct {
	DocumentNumber string  `json:"document_number" validate:"required"`
	CurrentBalance float64 `json:"current_balance" validate:"required,gt=0"`
	UserId         string  `json:"user_id"  validate:"required,uuid"`
}

// createAccount handles creating an account
func (h *Handler) createAccount() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		requestBody := &CreateAccountRequestData{}
		if ok := h.reader.ReadJSONAndValidate(w, r, requestBody); !ok {
			return
		}

		// Validate that the user exists in the database
		if !h.validateUserExists(ctx, w, requestBody.UserId) {
			return
		}

		// Create the account and respond
		h.createAndRespondAccount(ctx, w, requestBody)
	}
}

// createAndRespondAccount creates the account in the database and sends the response
func (h *Handler) createAndRespondAccount(ctx context.Context, w http.ResponseWriter, requestBody *CreateAccountRequestData) {
	accountDetails, err := h.repository.createAccount(ctx, models.CreateAccountParams{
		DocumentNumber: requestBody.DocumentNumber,
		CurrentBalance: requestBody.CurrentBalance,
		UserID:         requestBody.UserId,
	})

	if err != nil {
		log.Printf("createAndRespondAccount: failed to create an account: %v", err)
		h.writer.Internal(w, &response.APIError{
			Code:    response.DefaultErrorCode,
			Message: "Failed to create account.",
		})
		return
	}

	h.writer.Ok(w, accountDetails)
}

// validateUserExists checks if the provided user ID exists in the database
func (h *Handler) validateUserExists(ctx context.Context, w http.ResponseWriter, userID string) bool {
	userExists, err := h.repository.userExists(ctx, userID)
	if err != nil {
		log.Printf("validateUserExists: failed to check user existence: %v", err)
		h.writer.Internal(w, &response.APIError{
			Code:    response.DefaultErrorCode,
			Message: "Failed to validate user ID.",
		})
		return false
	}

	if !userExists {
		log.Printf("validateUserExists: user %s does not exist", userID)
		h.writer.NotFound(w, &response.APIError{
			Code:    response.ErrUserNotFound,
			Message: errUserNotFound.Error(),
		})
		return false
	}

	return true
}
