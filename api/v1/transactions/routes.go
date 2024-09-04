package transactions

import (
	"net/http"

	"github.com/gorilla/mux"
)

func Routes(r *mux.Router, h *Handler) {
	r.HandleFunc("", h.createTransaction()).Methods(http.MethodPost)
	r.HandleFunc("/{transactionID}", h.getTransactionDetails()).Methods(http.MethodGet)
}
