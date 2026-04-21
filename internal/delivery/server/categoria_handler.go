package server

import (
	"net/http"
	"strconv"

	"finance-management/internal/delivery/requests"
	"finance-management/tools/helpers"

	"github.com/gorilla/mux"
)

func (h *HTTPServer) GetCategoriasByUsuarioID(w http.ResponseWriter, r *http.Request) {
	usuarioID, err := strconv.Atoi(mux.Vars(r)["usuario_id"])
	if err != nil || usuarioID <= 0 {
		writeErrorResponse(w, http.StatusBadRequest, "invalid usuario_id", nil)
		return
	}

	categorias, err := h.categoriaService.GetAllByUsuarioID(r.Context(), usuarioID)
	if err != nil {
		h.logger.Error("could not get categorias", helpers.ErrLoggingKey, err)
		writeErrorResponse(w, http.StatusBadRequest, "could not get categorias", err.Error())
		return
	}

	writeSuccessResponse(w, http.StatusOK, categorias)
}

func (h *HTTPServer) CreateCategoria(w http.ResponseWriter, r *http.Request) {
	var request requests.CreateCategoryRequest
	if err := decodeJSON(r, &request); err != nil {
		writeErrorResponse(w, http.StatusBadRequest, "invalid request body", err.Error())
		return
	}

	if err := h.validator.Struct(request); err != nil {
		writeErrorResponse(w, http.StatusBadRequest, "invalid request data", helpers.GetErrorValidations(err))
		return
	}

	categoria, err := h.categoriaService.Create(r.Context(), request)
	if err != nil {
		h.logger.Error("could not create categoria", helpers.ErrLoggingKey, err)
		writeErrorResponse(w, http.StatusBadRequest, "could not create categoria", err.Error())
		return
	}

	writeSuccessResponse(w, http.StatusCreated, categoria)
}

func (h *HTTPServer) UpdateCategoria(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil || id <= 0 {
		writeErrorResponse(w, http.StatusBadRequest, "invalid categoria id", nil)
		return
	}

	var request requests.UpdateCategoryRequest
	if err := decodeJSON(r, &request); err != nil {
		writeErrorResponse(w, http.StatusBadRequest, "invalid request body", err.Error())
		return
	}
	request.ID = id

	if err := h.validator.Struct(request); err != nil {
		writeErrorResponse(w, http.StatusBadRequest, "invalid request data", helpers.GetErrorValidations(err))
		return
	}

	categoria, err := h.categoriaService.Update(r.Context(), request)
	if err != nil {
		h.logger.Error("could not update categoria", helpers.ErrLoggingKey, err)
		writeErrorResponse(w, http.StatusBadRequest, "could not update categoria", err.Error())
		return
	}

	writeSuccessResponse(w, http.StatusOK, categoria)
}

func (h *HTTPServer) DeleteCategoria(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil || id <= 0 {
		writeErrorResponse(w, http.StatusBadRequest, "invalid categoria id", nil)
		return
	}

	if err := h.categoriaService.Delete(r.Context(), id); err != nil {
		h.logger.Error("could not delete categoria", helpers.ErrLoggingKey, err)
		writeErrorResponse(w, http.StatusBadRequest, "could not delete categoria", err.Error())
		return
	}

	writeSuccessResponse(w, http.StatusOK, map[string]string{"message": "categoria deleted successfully"})
}
