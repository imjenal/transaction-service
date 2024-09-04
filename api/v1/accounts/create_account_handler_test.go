package accounts

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
	dummyUserId    = "88e0e837-e7f2-47b1-a08c-3af267c03088"
	dummyAccountId = "115be6d7-6d9a-4391-b3ee-1d753ac7d611"
)

func TestCreateAccountHandler_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock.NewMockQuerier(ctrl)
	writer := response.NewJSONWriter()
	v := validator.New()
	reader := request.NewReader(writer, v)
	handler := NewHandler(reader, writer, &Repository{querier: mockRepo})

	// Prepare mock responses
	mockRepo.EXPECT().UserExists(gomock.Any(), dummyUserId).Return(true, nil)
	mockRepo.EXPECT().CreateAccount(gomock.Any(), gomock.Any()).Return(&models.Account{Uuid: dummyAccountId}, nil)

	// Prepare the request
	requestBody, _ := json.Marshal(CreateAccountRequestData{
		DocumentNumber: "1234567890",
		CurrentBalance: 1000.0,
		UserId:         dummyUserId,
	})

	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(requestBody))
	rr := httptest.NewRecorder()

	// Call the handler
	handler.createAccount()(rr, req)

	// Check the results
	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Contains(t, rr.Body.String(), dummyAccountId)
}

func TestCreateAccountHandler_ValidationError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock.NewMockQuerier(ctrl)
	writer := response.NewJSONWriter()
	v := validator.New()
	reader := request.NewReader(writer, v)
	handler := NewHandler(reader, writer, &Repository{querier: mockRepo})

	// Prepare the invalid request
	requestBody, _ := json.Marshal(CreateAccountRequestData{
		DocumentNumber: "Doc131",
		CurrentBalance: 100.0,
		UserId:         "invalid-user-id",
	})

	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(requestBody))
	rr := httptest.NewRecorder()

	// Call the handler
	handler.createAccount()(rr, req)

	// Check the results
	assert.Equal(t, http.StatusUnprocessableEntity, rr.Code)
}

func TestCreateAccountHandler_UserNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock.NewMockQuerier(ctrl)
	writer := response.NewJSONWriter()
	v := validator.New()
	reader := request.NewReader(writer, v)
	handler := NewHandler(reader, writer, &Repository{querier: mockRepo})

	// Prepare mock response
	mockRepo.EXPECT().UserExists(gomock.Any(), dummyUserId).Return(false, nil)

	// Prepare the request
	requestBody, _ := json.Marshal(CreateAccountRequestData{
		DocumentNumber: "1234567890",
		CurrentBalance: 1000.0,
		UserId:         dummyUserId,
	})

	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(requestBody))
	rr := httptest.NewRecorder()

	// Call the handler
	handler.createAccount()(rr, req)

	// Check the results
	assert.Equal(t, http.StatusNotFound, rr.Code)
}

func TestCreateAccountHandler_DBError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock.NewMockQuerier(ctrl)
	writer := response.NewJSONWriter()
	v := validator.New()
	reader := request.NewReader(writer, v)
	handler := NewHandler(reader, writer, &Repository{querier: mockRepo})

	// Mock database error
	mockRepo.EXPECT().UserExists(gomock.Any(), dummyUserId).Return(true, nil)
	mockRepo.EXPECT().CreateAccount(gomock.Any(), gomock.Any()).Return(nil, errors.New("database error"))

	// Prepare the request
	requestBody, _ := json.Marshal(CreateAccountRequestData{
		DocumentNumber: "1234567890",
		CurrentBalance: 1000.0,
		UserId:         dummyUserId,
	})

	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(requestBody))
	rr := httptest.NewRecorder()

	// Call the handler
	handler.createAccount()(rr, req)

	// Check the results
	assert.Equal(t, http.StatusInternalServerError, rr.Code)
	assert.Contains(t, rr.Body.String(), "Failed to create account.")
}
