package accounts

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

const (
	dummyAccountID = "115be6d7-6d9a-4391-b3ee-1d753ac7d611"
)

func TestGetAccountHandler_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock.NewMockQuerier(ctrl)
	writer := response.NewJSONWriter()
	v := validator.New()
	reader := request.NewReader(writer, v)
	handler := NewHandler(reader, writer, &Repository{querier: mockRepo})

	// Prepare mock response
	mockRepo.EXPECT().GetAccountDetailsByUUID(gomock.Any(), dummyAccountID).Return(&models.Account{Uuid: dummyAccountID}, nil)

	// Prepare the request
	req := httptest.NewRequest(http.MethodGet, "/accounts/"+dummyAccountID, nil)
	rr := httptest.NewRecorder()

	// Set mux variables
	req = mux.SetURLVars(req, map[string]string{"accountID": dummyAccountID})

	// Call the handler
	handler.getAccountDetails()(rr, req)

	// Check the results
	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Contains(t, rr.Body.String(), dummyAccountID)
}

func TestGetAccountHandler_AccountNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock.NewMockQuerier(ctrl)
	writer := response.NewJSONWriter()
	v := validator.New()
	reader := request.NewReader(writer, v)
	handler := NewHandler(reader, writer, &Repository{querier: mockRepo})

	// Prepare mock response
	mockRepo.EXPECT().GetAccountDetailsByUUID(gomock.Any(), dummyAccountID).Return(nil, errAccountNotFound)

	// Prepare the request
	req := httptest.NewRequest(http.MethodGet, "/accounts/"+dummyAccountID, nil)
	rr := httptest.NewRecorder()

	// Set mux variables
	req = mux.SetURLVars(req, map[string]string{"accountID": dummyAccountID})

	// Call the handler
	handler.getAccountDetails()(rr, req)

	// Check the results
	assert.Equal(t, http.StatusNotFound, rr.Code)
}

func TestGetAccountHandler_DBError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock.NewMockQuerier(ctrl)
	writer := response.NewJSONWriter()
	v := validator.New()
	reader := request.NewReader(writer, v)
	handler := NewHandler(reader, writer, &Repository{querier: mockRepo})

	// Prepare mock response for database error
	mockRepo.EXPECT().GetAccountDetailsByUUID(gomock.Any(), dummyAccountID).Return(nil, errors.New("database error"))

	// Prepare the request
	req := httptest.NewRequest(http.MethodGet, "/accounts/"+dummyAccountID, nil)
	rr := httptest.NewRecorder()

	// Set mux variables
	req = mux.SetURLVars(req, map[string]string{"accountID": dummyAccountID})

	// Call the handler
	handler.getAccountDetails()(rr, req)

	// Check the results
	assert.Equal(t, http.StatusInternalServerError, rr.Code)
	assert.Contains(t, rr.Body.String(), "Failed to fetch account details.")
}
