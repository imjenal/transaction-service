package accounts

import (
	"github.com/imjenal/transaction-service/pkg/http/request"
	"github.com/imjenal/transaction-service/pkg/http/response"
)

type Handler struct {
	reader     *request.Reader
	writer     *response.JSONWriter
	repository *Repository
}

func NewHandler(reader *request.Reader, writer *response.JSONWriter, repository *Repository) *Handler {
	return &Handler{
		reader:     reader,
		writer:     writer,
		repository: repository,
	}
}
