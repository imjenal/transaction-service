package api

import (
	"github.com/imjenal/transaction-service/pkg/validator"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/imjenal/transaction-service/api/v1/accounts"
	"github.com/imjenal/transaction-service/api/v1/transactions"
	"github.com/imjenal/transaction-service/internal/app"
	"github.com/imjenal/transaction-service/internal/db"
	"github.com/imjenal/transaction-service/internal/db/models"
	"github.com/imjenal/transaction-service/pkg/http/request"
	"github.com/imjenal/transaction-service/pkg/http/response"
)

type Params struct {
	DB        *db.DB
	Reader    *request.Reader
	Writer    *response.JSONWriter
	Validator *validator.Validator
}

func Routes(r *mux.Router, params *Params) {
	querier := models.New(params.DB.Conn)

	// Add the API health check route at the top level
	r.HandleFunc("/health", healthCheck(params.Writer))

	// Create a /v1/ sub-router for the API
	v1Router := r.PathPrefix("/v1/").Subrouter()

	pathValidatorMiddleware := validator.NewPathValidator(params.Validator, params.Writer, map[string]string{
		"transactionID": "uuid4",
		"accountID":     "uuid4",
	})
	v1Router.Use(pathValidatorMiddleware)

	// All repositories are initialized here
	accountsRepo := accounts.NewRepository(querier)
	transactionsRepo := transactions.NewRepository(querier)

	// All handlers are initialized here
	accountsHandler := accounts.NewHandler(params.Reader, params.Writer, accountsRepo)
	transactionsHandler := transactions.NewHandler(params.Reader, params.Writer, transactionsRepo)

	// All routes are added here
	accounts.Routes(v1Router.PathPrefix("/accounts").Subrouter(), accountsHandler)
	transactions.Routes(v1Router.PathPrefix("/transactions").Subrouter(), transactionsHandler)

}

// healthCheck returns a handler that returns a 200 OK response with version and commit hash
func healthCheck(jsonWriter *response.JSONWriter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		jsonWriter.Ok(w, map[string]any{"version": app.Version()})
	}
}
