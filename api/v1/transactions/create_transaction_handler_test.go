package transactions

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/imjenal/transaction-service/internal/db/models"
	"github.com/imjenal/transaction-service/internal/db/models/mock"
	"github.com/imjenal/transaction-service/pkg/http/request"
	"github.com/imjenal/transaction-service/pkg/http/response"
	"github.com/imjenal/transaction-service/pkg/validator"
	"github.com/stretchr/testify/assert"
)

const (
	dummyAccountId     = "115be6d7-6d9a-4391-b3ee-1d753ac7d611"
	dummyOperationType = int64(1)
	dummyTransactionID = "98a0f8e7-6e28-4d4f-872b-4d28b3d5ee66"
)

func TestCreateTransactionHandler_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock.NewMockQuerier(ctrl)
	writer := response.NewJSONWriter()
	v := validator.New()
	reader := request.NewReader(writer, v)
	handler := NewHandler(reader, writer, &Repository{querier: mockRepo})

	// Prepare mock responses
	mockRepo.EXPECT().AccountExists(gomock.Any(), dummyAccountId).Return(true, nil)
	mockRepo.EXPECT().GetOperationTypeAmountBehavior(gomock.Any(), dummyOperationType).Return(models.AmountBehaviorNEGATIVE, nil)
	mockRepo.EXPECT().CreateTransaction(gomock.Any(), gomock.Any()).Return(&models.Transaction{Uuid: dummyTransactionID}, nil)

	// Prepare the request
	requestBody, _ := json.Marshal(CreateTransactionRequestData{
		AccountId:       dummyAccountId,
		OperationTypeId: dummyOperationType,
		Amount:          100.0,
	})

	req := httptest.NewRequest(http.MethodPost, "/transactions", bytes.NewReader(requestBody))
	rr := httptest.NewRecorder()

	// Call the handler
	handler.createTransaction()(rr, req)

	// Check the results
	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Contains(t, rr.Body.String(), dummyTransactionID)
}

func TestCreateTransactionHandler_ValidationError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock.NewMockQuerier(ctrl)
	writer := response.NewJSONWriter()
	v := validator.New()
	reader := request.NewReader(writer, v)
	handler := NewHandler(reader, writer, &Repository{querier: mockRepo})

	// Prepare the invalid request
	requestBody, _ := json.Marshal(CreateTransactionRequestData{
		AccountId:       "", // Invalid account ID
		OperationTypeId: dummyOperationType,
		Amount:          -100.0, // Invalid amount (negative value)
	})

	req := httptest.NewRequest(http.MethodPost, "/transactions", bytes.NewReader(requestBody))
	rr := httptest.NewRecorder()

	// Call the handler
	handler.createTransaction()(rr, req)

	// Check the results
	assert.Equal(t, http.StatusUnprocessableEntity, rr.Code)
}

func TestCreateTransactionHandler_AccountNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock.NewMockQuerier(ctrl)
	writer := response.NewJSONWriter()
	v := validator.New()
	reader := request.NewReader(writer, v)
	handler := NewHandler(reader, writer, &Repository{querier: mockRepo})

	// Prepare mock response
	mockRepo.EXPECT().AccountExists(gomock.Any(), dummyAccountId).Return(false, nil)

	// Prepare the request
	requestBody, _ := json.Marshal(CreateTransactionRequestData{
		AccountId:       dummyAccountId,
		OperationTypeId: dummyOperationType,
		Amount:          100.0,
	})

	req := httptest.NewRequest(http.MethodPost, "/transactions", bytes.NewReader(requestBody))
	rr := httptest.NewRecorder()

	// Call the handler
	handler.createTransaction()(rr, req)

	// Check the results
	assert.Equal(t, http.StatusNotFound, rr.Code)
}

func TestCreateTransactionHandler_OperationTypeNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock.NewMockQuerier(ctrl)
	writer := response.NewJSONWriter()
	v := validator.New()
	reader := request.NewReader(writer, v)
	handler := NewHandler(reader, writer, &Repository{querier: mockRepo})

	// Prepare mock responses
	mockRepo.EXPECT().AccountExists(gomock.Any(), dummyAccountId).Return(true, nil)
	mockRepo.EXPECT().GetOperationTypeAmountBehavior(gomock.Any(), dummyOperationType).Return(models.AmountBehavior(""), errOperationTypeNotFound)

	// Prepare the request
	requestBody, _ := json.Marshal(CreateTransactionRequestData{
		AccountId:       dummyAccountId,
		OperationTypeId: dummyOperationType,
		Amount:          100.0,
	})

	req := httptest.NewRequest(http.MethodPost, "/transactions", bytes.NewReader(requestBody))
	rr := httptest.NewRecorder()

	// Call the handler
	handler.createTransaction()(rr, req)

	// Check the results
	assert.Equal(t, http.StatusNotFound, rr.Code)
}

func TestCreateTransactionHandler_DBError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock.NewMockQuerier(ctrl)
	writer := response.NewJSONWriter()
	v := validator.New()
	reader := request.NewReader(writer, v)
	handler := NewHandler(reader, writer, &Repository{querier: mockRepo})

	// Mock database error during account validation
	mockRepo.EXPECT().AccountExists(gomock.Any(), dummyAccountId).Return(true, nil)
	mockRepo.EXPECT().GetOperationTypeAmountBehavior(gomock.Any(), dummyOperationType).Return(models.AmountBehaviorNEGATIVE, nil)
	mockRepo.EXPECT().CreateTransaction(gomock.Any(), gomock.Any()).Return(nil, errors.New("database error"))

	// Prepare the request
	requestBody, _ := json.Marshal(CreateTransactionRequestData{
		AccountId:       dummyAccountId,
		OperationTypeId: dummyOperationType,
		Amount:          100.0,
	})

	req := httptest.NewRequest(http.MethodPost, "/transactions", bytes.NewReader(requestBody))
	rr := httptest.NewRecorder()

	// Call the handler
	handler.createTransaction()(rr, req)

	// Check the results
	assert.Equal(t, http.StatusInternalServerError, rr.Code)
	assert.Contains(t, rr.Body.String(), "Failed to create transaction.")
}
