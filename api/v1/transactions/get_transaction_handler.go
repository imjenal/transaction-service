package transactions

import (
	"context"
	"errors"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/imjenal/transaction-service/pkg/http/response"
)

// getTransactionDetails handles fetching the transaction details
func (h *Handler) getTransactionDetails() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		transactionID := mux.Vars(r)["transactionID"]

		// Fetch and respond with the transaction details
		h.fetchAndRespondTransactionDetails(r.Context(), w, transactionID)

	}
}

// fetchAndRespondTransactionDetails fetches transaction details from the repository and responds to the client
func (h *Handler) fetchAndRespondTransactionDetails(ctx context.Context, w http.ResponseWriter, transactionID string) {
	txnDetails, err := h.repository.getTransactionDetails(ctx, transactionID)
	if errors.Is(err, errTransactionNotFound) {
		log.Printf("fetchAndRespondTransactionDetails: transaction %s not found", transactionID)
		h.writer.NotFound(w, &response.APIError{
			Code:    response.ErrTransactionNotFound,
			Message: errTransactionNotFound.Error(),
		})
		return
	}

	if err != nil {
		log.Printf("fetchAndRespondTransactionDetails: failed to fetch transaction details: %v", err)
		h.writer.Internal(w, &response.APIError{
			Code:    response.DefaultErrorCode,
			Message: "Failed to fetch transaction details.",
		})
		return
	}

	h.writer.Ok(w, txnDetails)
}
