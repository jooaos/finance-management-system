package server

import (
	"net/http"

	"finance-management/internal/delivery/requests"
	"finance-management/tools/helpers"
)

func (h *HTTPServer) CreateTransacao(w http.ResponseWriter, r *http.Request) {
	var request requests.CreateTransactionRequest
	if err := decodeJSON(r, &request); err != nil {
		writeErrorResponse(w, http.StatusBadRequest, "invalid request body", err.Error())
		return
	}

	if err := h.validator.Struct(request); err != nil {
		writeErrorResponse(w, http.StatusBadRequest, "invalid request data", helpers.GetErrorValidations(err))
		return
	}

	transacoes, err := h.transacaoService.Create(r.Context(), request)
	if err != nil {
		h.logger.Error("could not create transacao", helpers.ErrLoggingKey, err)
		writeErrorResponse(w, http.StatusBadRequest, "could not create transacao", err.Error())
		return
	}

	writeSuccessResponse(w, http.StatusCreated, transacoes)
}
