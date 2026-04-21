package server

import (
	"net/http"

	"finance-management/internal/delivery/requests"
	"finance-management/tools/helpers"
)

func (h *HTTPServer) CreateUsuario(w http.ResponseWriter, r *http.Request) {
	var request requests.CreateUserRequest
	if err := decodeJSON(r, &request); err != nil {
		writeErrorResponse(w, http.StatusBadRequest, "invalid request body", err.Error())
		return
	}

	if err := h.validator.Struct(request); err != nil {
		writeErrorResponse(w, http.StatusBadRequest, "invalid request data", helpers.GetErrorValidations(err))
		return
	}

	usuario, err := h.usuarioService.Create(r.Context(), request)
	if err != nil {
		h.logger.Error("could not create usuario", helpers.ErrLoggingKey, err)
		writeErrorResponse(w, http.StatusBadRequest, "could not create usuario", err.Error())
		return
	}

	writeSuccessResponse(w, http.StatusCreated, usuario)
}
