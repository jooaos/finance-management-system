package server

import (
	"log/slog"
	"net/http"

	"finance-management/internal/service"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
)

type HTTPServer struct {
	logger           *slog.Logger
	validator        *validator.Validate
	usuarioService   *service.UsuarioService
	transacaoService *service.TransacaoService
	categoriaService *service.CategoriaService
	orcamentoService *service.OrcamentoService
}

func NewHTTPServer(
	logger *slog.Logger,
	usuarioService *service.UsuarioService,
	transacaoService *service.TransacaoService,
	categoriaService *service.CategoriaService,
	orcamentoService *service.OrcamentoService,
) *HTTPServer {
	return &HTTPServer{
		logger:           logger,
		validator:        validator.New(),
		usuarioService:   usuarioService,
		transacaoService: transacaoService,
		categoriaService: categoriaService,
		orcamentoService: orcamentoService,
	}
}

func (h *HTTPServer) InitServer() *mux.Router {
	r := mux.NewRouter()
	s := r.PathPrefix("/api").Subrouter()

	s.HandleFunc("/usuarios", h.CreateUsuario).Methods("POST")
	s.HandleFunc("/transacoes", h.CreateTransacao).Methods("POST")
	s.HandleFunc("/usuarios/{usuario_id}/categorias", h.GetCategoriasByUsuarioID).Methods("GET")
	s.HandleFunc("/categorias", h.CreateCategoria).Methods("POST")
	s.HandleFunc("/categorias/{id}", h.DeleteCategoria).Methods("DELETE")
	s.HandleFunc("/categorias/{id}", h.UpdateCategoria).Methods("PUT")
	s.HandleFunc("/orcamentos", h.CreateOrcamento).Methods("POST")
	s.HandleFunc("/usuarios/{usuario_id}/orcamentos", h.GetOrcamentosByUsuarioIDAndMes).Methods("GET")
	s.HandleFunc("/usuarios/{usuario_id}/orcamentos/total", h.GetTotalOrcamentoByUsuarioIDAndMes).Methods("GET")
	s.HandleFunc("/orcamentos/{id}", h.UpdateOrcamento).Methods("PUT")
	s.HandleFunc("/usuarios/{usuario_id}/relatorios/mensal", h.GetMonthlyTransactionSummary).Methods("GET")
	s.HandleFunc("/usuarios/{usuario_id}/relatorios/gastos", h.GetMonthlySpendingProgress).Methods("GET")
	s.HandleFunc("/usuarios/{usuario_id}/relatorios/categorias", h.GetCategoryMonthlySummary).Methods("GET")
	s.HandleFunc("/usuarios/{usuario_id}/projecao/comprometimento", h.GetCommitmentProjection).Methods("GET")

	r.PathPrefix("/").Handler(http.FileServer(http.Dir("frontend")))

	return r
}
