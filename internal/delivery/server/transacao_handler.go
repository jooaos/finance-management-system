package server

import (
	"net/http"
	"strconv"

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

func (h *HTTPServer) GetMonthlyTransactionSummary(w http.ResponseWriter, r *http.Request) {
	usuarioID, mes, ok := h.getUsuarioIDAndMesFromRequest(w, r)
	if !ok {
		return
	}

	summary, err := h.transacaoService.GetMonthlySummary(r.Context(), usuarioID, mes)
	if err != nil {
		h.logger.Error("could not get monthly transaction summary", helpers.ErrLoggingKey, err)
		writeErrorResponse(w, http.StatusBadRequest, "could not get monthly transaction summary", err.Error())
		return
	}

	writeSuccessResponse(w, http.StatusOK, summary)
}

func (h *HTTPServer) GetMonthlySpendingProgress(w http.ResponseWriter, r *http.Request) {
	usuarioID, mes, ok := h.getUsuarioIDAndMesFromRequest(w, r)
	if !ok {
		return
	}

	progress, err := h.transacaoService.GetMonthlySpendingProgress(r.Context(), usuarioID, mes)
	if err != nil {
		h.logger.Error("could not get monthly spending progress", helpers.ErrLoggingKey, err)
		writeErrorResponse(w, http.StatusBadRequest, "could not get monthly spending progress", err.Error())
		return
	}

	writeSuccessResponse(w, http.StatusOK, progress)
}

func (h *HTTPServer) GetCategoryMonthlySummary(w http.ResponseWriter, r *http.Request) {
	usuarioID, mes, ok := h.getUsuarioIDAndMesFromRequest(w, r)
	if !ok {
		return
	}

	summary, err := h.transacaoService.GetCategoryMonthlySummary(r.Context(), usuarioID, mes)
	if err != nil {
		h.logger.Error("could not get category monthly summary", helpers.ErrLoggingKey, err)
		writeErrorResponse(w, http.StatusBadRequest, "could not get category monthly summary", err.Error())
		return
	}

	writeSuccessResponse(w, http.StatusOK, summary)
}

func (h *HTTPServer) GetCommitmentProjection(w http.ResponseWriter, r *http.Request) {
	usuarioID, mes, ok := h.getUsuarioIDAndMesFromRequest(w, r)
	if !ok {
		return
	}

	months := 4
	monthsParam := r.URL.Query().Get("meses")
	if monthsParam != "" {
		parsedMonths, err := strconv.Atoi(monthsParam)
		if err != nil || parsedMonths <= 0 {
			writeErrorResponse(w, http.StatusBadRequest, "invalid meses query param", nil)
			return
		}
		months = parsedMonths
	}

	projection, err := h.transacaoService.GetCommitmentProjection(r.Context(), usuarioID, mes, months)
	if err != nil {
		h.logger.Error("could not get commitment projection", helpers.ErrLoggingKey, err)
		writeErrorResponse(w, http.StatusBadRequest, "could not get commitment projection", err.Error())
		return
	}

	writeSuccessResponse(w, http.StatusOK, projection)
}
