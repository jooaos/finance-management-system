package service

import (
	"context"
	"fmt"
	"time"

	"finance-management/internal/delivery/requests"
	"finance-management/internal/domain"
	"finance-management/internal/repository"
)

type OrcamentoService struct {
	orcamentosRepository *repository.OrcamentosRepository
}

func NewOrcamentoService(orcamentosRepository *repository.OrcamentosRepository) *OrcamentoService {
	return &OrcamentoService{
		orcamentosRepository: orcamentosRepository,
	}
}

func (b *OrcamentoService) Create(ctx context.Context, request requests.CreateBudgetRequest) (*domain.Orcamento, error) {
	orcamento := &domain.Orcamento{
		UsuarioID:   request.UsuarioID,
		CategoriaID: request.CategoriaID,
		Limite:      request.Limite,
		Mes:         normalizeMes(request.Mes),
	}

	createdOrcamento, err := b.orcamentosRepository.Create(ctx, orcamento)
	if err != nil {
		return nil, fmt.Errorf("failed to create budget: %w", err)
	}

	return createdOrcamento, nil
}

func (b *OrcamentoService) GetByUsuarioIDAndMes(ctx context.Context, usuarioID int, mes time.Time) ([]domain.Orcamento, error) {
	orcamentos, err := b.orcamentosRepository.GetByUsuarioIDAndMes(ctx, usuarioID, normalizeMes(mes))
	if err != nil {
		return nil, fmt.Errorf("failed to get budgets by user id and month: %w", err)
	}

	return orcamentos, nil
}

func (b *OrcamentoService) GetTotalPlanejadoByUsuarioIDAndMes(ctx context.Context, usuarioID int, mes time.Time) (float64, error) {
	total, err := b.orcamentosRepository.GetTotalPlanejadoByUsuarioIDAndMes(ctx, usuarioID, normalizeMes(mes))
	if err != nil {
		return 0, fmt.Errorf("failed to get total planned budget by month: %w", err)
	}

	return total, nil
}

func (b *OrcamentoService) Update(ctx context.Context, request requests.UpdateBudgetRequest) (*domain.Orcamento, error) {
	orcamento := &domain.Orcamento{
		ID:          request.ID,
		UsuarioID:   request.UsuarioID,
		CategoriaID: request.CategoriaID,
		Limite:      request.Limite,
		Mes:         normalizeMes(request.Mes),
	}

	updatedOrcamento, err := b.orcamentosRepository.Update(ctx, orcamento)
	if err != nil {
		return nil, fmt.Errorf("failed to update budget: %w", err)
	}

	return updatedOrcamento, nil
}

func (b *OrcamentoService) GetByID(ctx context.Context, id int) (*domain.Orcamento, error) {
	orcamento, err := b.orcamentosRepository.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get budget by id: %w", err)
	}

	return orcamento, nil
}

func normalizeMes(mes time.Time) time.Time {
	return time.Date(mes.Year(), mes.Month(), 1, 0, 0, 0, 0, mes.Location())
}
