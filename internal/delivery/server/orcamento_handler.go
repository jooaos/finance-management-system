package server

import (
	"net/http"
	"strconv"
	"time"

	"finance-management/internal/delivery/requests"
	"finance-management/tools/helpers"

	"github.com/gorilla/mux"
)

const queryDateLayout = "2006-01-02"

func (h *HTTPServer) CreateOrcamento(w http.ResponseWriter, r *http.Request) {
	var request requests.CreateBudgetRequest
	if err := decodeJSON(r, &request); err != nil {
		writeErrorResponse(w, http.StatusBadRequest, "invalid request body", err.Error())
		return
	}

	if err := h.validator.Struct(request); err != nil {
		writeErrorResponse(w, http.StatusBadRequest, "invalid request data", helpers.GetErrorValidations(err))
		return
	}

	orcamento, err := h.orcamentoService.Create(r.Context(), request)
	if err != nil {
		h.logger.Error("could not create orcamento", helpers.ErrLoggingKey, err)
		writeErrorResponse(w, http.StatusBadRequest, "could not create orcamento", err.Error())
		return
	}

	writeSuccessResponse(w, http.StatusCreated, orcamento)
}

func (h *HTTPServer) GetOrcamentosByUsuarioIDAndMes(w http.ResponseWriter, r *http.Request) {
	usuarioID, mes, ok := h.getUsuarioIDAndMesFromRequest(w, r)
	if !ok {
		return
	}

	orcamentos, err := h.orcamentoService.GetByUsuarioIDAndMes(r.Context(), usuarioID, mes)
	if err != nil {
		h.logger.Error("could not get orcamentos", helpers.ErrLoggingKey, err)
		writeErrorResponse(w, http.StatusBadRequest, "could not get orcamentos", err.Error())
		return
	}

	writeSuccessResponse(w, http.StatusOK, orcamentos)
}

func (h *HTTPServer) GetTotalOrcamentoByUsuarioIDAndMes(w http.ResponseWriter, r *http.Request) {
	usuarioID, mes, ok := h.getUsuarioIDAndMesFromRequest(w, r)
	if !ok {
		return
	}

	total, err := h.orcamentoService.GetTotalPlanejadoByUsuarioIDAndMes(r.Context(), usuarioID, mes)
	if err != nil {
		h.logger.Error("could not get total orcamento", helpers.ErrLoggingKey, err)
		writeErrorResponse(w, http.StatusBadRequest, "could not get total orcamento", err.Error())
		return
	}

	writeSuccessResponse(w, http.StatusOK, map[string]float64{"total_planejado": total})
}

func (h *HTTPServer) UpdateOrcamento(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil || id <= 0 {
		writeErrorResponse(w, http.StatusBadRequest, "invalid orcamento id", nil)
		return
	}

	var request requests.UpdateBudgetRequest
	if err := decodeJSON(r, &request); err != nil {
		writeErrorResponse(w, http.StatusBadRequest, "invalid request body", err.Error())
		return
	}
	request.ID = id

	if err := h.validator.Struct(request); err != nil {
		writeErrorResponse(w, http.StatusBadRequest, "invalid request data", helpers.GetErrorValidations(err))
		return
	}

	orcamento, err := h.orcamentoService.Update(r.Context(), request)
	if err != nil {
		h.logger.Error("could not update orcamento", helpers.ErrLoggingKey, err)
		writeErrorResponse(w, http.StatusBadRequest, "could not update orcamento", err.Error())
		return
	}

	writeSuccessResponse(w, http.StatusOK, orcamento)
}

func (h *HTTPServer) getUsuarioIDAndMesFromRequest(w http.ResponseWriter, r *http.Request) (int, time.Time, bool) {
	usuarioID, err := strconv.Atoi(mux.Vars(r)["usuario_id"])
	if err != nil || usuarioID <= 0 {
		writeErrorResponse(w, http.StatusBadRequest, "invalid usuario_id", nil)
		return 0, time.Time{}, false
	}

	mesParam := r.URL.Query().Get("mes")
	if mesParam == "" {
		writeErrorResponse(w, http.StatusBadRequest, "missing mes query param", nil)
		return 0, time.Time{}, false
	}

	mes, err := time.Parse(queryDateLayout, mesParam)
	if err != nil {
		writeErrorResponse(w, http.StatusBadRequest, "invalid mes query param, use YYYY-MM-DD", err.Error())
		return 0, time.Time{}, false
	}

	return usuarioID, mes, true
}
