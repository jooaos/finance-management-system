package repository

import (
	"context"
	"fmt"
	"time"

	"finance-management/internal/domain"

	"gorm.io/gorm"
)

var (
	orcamentosTable = "orcamentos"
)

type OrcamentosRepository struct {
	conn  *gorm.DB
	table string
}

func NewOrcamentosRepository(
	conn *gorm.DB,
) *OrcamentosRepository {
	return &OrcamentosRepository{
		conn:  conn,
		table: orcamentosTable,
	}
}

func (o *OrcamentosRepository) Create(ctx context.Context, orcamento *domain.Orcamento) (*domain.Orcamento, error) {
	if err := o.conn.WithContext(ctx).Table(o.table).Create(orcamento).Error; err != nil {
		return nil, fmt.Errorf("failed to create budget: %w", err)
	}

	return orcamento, nil
}

func (o *OrcamentosRepository) GetByUsuarioIDAndMes(ctx context.Context, usuarioID int, mes time.Time) ([]domain.Orcamento, error) {
	var orcamentos []domain.Orcamento
	if err := o.conn.WithContext(ctx).
		Table(o.table).
		Where("usuario_id = ? AND mes = ?", usuarioID, mes).
		Find(&orcamentos).Error; err != nil {
		return nil, fmt.Errorf("failed to get budgets by user id and month: %w", err)
	}

	return orcamentos, nil
}

func (o *OrcamentosRepository) GetTotalPlanejadoByUsuarioIDAndMes(ctx context.Context, usuarioID int, mes time.Time) (float64, error) {
	var total float64
	if err := o.conn.WithContext(ctx).
		Table(o.table).
		Select("COALESCE(SUM(limite), 0)").
		Where("usuario_id = ? AND mes = ?", usuarioID, mes).
		Scan(&total).Error; err != nil {
		return 0, fmt.Errorf("failed to get total planned budget by month: %w", err)
	}

	return total, nil
}

func (o *OrcamentosRepository) Update(ctx context.Context, orcamento *domain.Orcamento) (*domain.Orcamento, error) {
	updates := map[string]interface{}{
		"usuario_id":   orcamento.UsuarioID,
		"categoria_id": orcamento.CategoriaID,
		"limite":       orcamento.Limite,
		"mes":          orcamento.Mes,
	}

	if err := o.conn.WithContext(ctx).
		Table(o.table).
		Where("id = ?", orcamento.ID).
		Updates(updates).Error; err != nil {
		return nil, fmt.Errorf("failed to update budget: %w", err)
	}

	return o.GetByID(ctx, orcamento.ID)
}

func (o *OrcamentosRepository) GetByID(ctx context.Context, id int) (*domain.Orcamento, error) {
	var orcamento domain.Orcamento
	if err := o.conn.WithContext(ctx).Table(o.table).Where("id = ?", id).First(&orcamento).Error; err != nil {
		return nil, fmt.Errorf("failed to get budget by id: %w", err)
	}

	return &orcamento, nil
}
