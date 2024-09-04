package accounts

import (
	"net/http"

	"github.com/gorilla/mux"
)

func Routes(r *mux.Router, h *Handler) {
	r.HandleFunc("/{accountID}", h.getAccountDetails()).Methods(http.MethodGet)
	r.HandleFunc("", h.createAccount()).Methods(http.MethodPost)
}
