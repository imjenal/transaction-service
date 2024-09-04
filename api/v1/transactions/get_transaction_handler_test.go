package transactions

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/imjenal/transaction-service/internal/db/models"
	"github.com/imjenal/transaction-service/internal/db/models/mock"
	"github.com/imjenal/transaction-service/pkg/http/request"
	"github.com/imjenal/transaction-service/pkg/http/response"
	"github.com/imjenal/transaction-service/pkg/validator"
	"github.com/stretchr/testify/assert"
)

func TestGetTransactionHandler_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock.NewMockQuerier(ctrl)
	writer := response.NewJSONWriter()
	v := validator.New()
	reader := request.NewReader(writer, v)
	handler := NewHandler(reader, writer, &Repository{querier: mockRepo})

	// Prepare mock response
	mockRepo.EXPECT().GetTransactionDetailsByTransactionId(gomock.Any(), dummyTransactionID).Return(&models.Transaction{Uuid: dummyTransactionID}, nil)

	// Prepare the request
	req := httptest.NewRequest(http.MethodGet, "/transactions/"+dummyTransactionID, nil)
	rr := httptest.NewRecorder()

	// Set mux variables
	req = mux.SetURLVars(req, map[string]string{"transactionID": dummyTransactionID})

	// Call the handler
	handler.getTransactionDetails()(rr, req)

	// Check the results
	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Contains(t, rr.Body.String(), dummyTransactionID)
}

func TestGetTransactionHandler_TransactionNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock.NewMockQuerier(ctrl)
	writer := response.NewJSONWriter()
	v := validator.New()
	reader := request.NewReader(writer, v)
	handler := NewHandler(reader, writer, &Repository{querier: mockRepo})

	// Prepare mock response
	mockRepo.EXPECT().GetTransactionDetailsByTransactionId(gomock.Any(), dummyTransactionID).Return(nil, errTransactionNotFound)

	// Prepare the request
	req := httptest.NewRequest(http.MethodGet, "/transactions/"+dummyTransactionID, nil)
	rr := httptest.NewRecorder()

	// Set mux variables
	req = mux.SetURLVars(req, map[string]string{"transactionID": dummyTransactionID})

	// Call the handler
	handler.getTransactionDetails()(rr, req)

	// Check the results
	assert.Equal(t, http.StatusNotFound, rr.Code)
}

func TestGetTransactionHandler_DBError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock.NewMockQuerier(ctrl)
	writer := response.NewJSONWriter()
	v := validator.New()
	reader := request.NewReader(writer, v)
	handler := NewHandler(reader, writer, &Repository{querier: mockRepo})

	// Prepare mock response for database error
	mockRepo.EXPECT().GetTransactionDetailsByTransactionId(gomock.Any(), dummyTransactionID).Return(nil, errors.New("database error"))

	// Prepare the request
	req := httptest.NewRequest(http.MethodGet, "/transactions/"+dummyTransactionID, nil)
	rr := httptest.NewRecorder()

	// Set mux variables
	req = mux.SetURLVars(req, map[string]string{"transactionID": dummyTransactionID})

	// Call the handler
	handler.getTransactionDetails()(rr, req)

	// Check the results
	assert.Equal(t, http.StatusInternalServerError, rr.Code)
	assert.Contains(t, rr.Body.String(), "Failed to fetch transaction details.")
}
