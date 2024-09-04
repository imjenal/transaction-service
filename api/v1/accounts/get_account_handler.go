package accounts

import (
	"context"
	"errors"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/imjenal/transaction-service/pkg/http/response"
)

// getAccountDetails handles fetching the account details
func (h *Handler) getAccountDetails() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		accountID := mux.Vars(r)["accountID"]

		// Fetch and respond with the account details
		h.fetchAndRespondAccountDetails(r.Context(), w, accountID)
	}
}

// fetchAndRespondAccountDetails fetches the account details from the repository and responds to the client
func (h *Handler) fetchAndRespondAccountDetails(ctx context.Context, w http.ResponseWriter, accountID string) {
	accountDetails, err := h.repository.getAccountDetails(ctx, accountID)
	if errors.Is(err, errAccountNotFound) {
		log.Printf("fetchAndRespondAccountDetails: account %s not found", accountID)
		h.writer.NotFound(w, &response.APIError{
			Code:    response.ErrAccountNotFound,
			Message: errAccountNotFound.Error(),
		})
		return
	}

	if err != nil {
		log.Printf("fetchAndRespondAccountDetails: failed to fetch account details: %v", err)
		h.writer.Internal(w, &response.APIError{
			Code:    response.DefaultErrorCode,
			Message: "Failed to fetch account details.",
		})
		return
	}

	h.writer.Ok(w, accountDetails)
}
