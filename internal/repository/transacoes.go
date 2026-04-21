package repository

import (
	"context"
	"fmt"
	"time"

	"finance-management/internal/domain"

	"gorm.io/gorm"
)

var (
	transacoesTable = "transacoes"
)

type TransacoesRepository struct {
	conn  *gorm.DB
	table string
}

func NewTransacoesRepository(
	conn *gorm.DB,
) *TransacoesRepository {
	return &TransacoesRepository{
		conn:  conn,
		table: transacoesTable,
	}
}

func (t *TransacoesRepository) Create(ctx context.Context, transacao *domain.Transacao) (*domain.Transacao, error) {
	if err := t.conn.WithContext(ctx).Table(t.table).Create(transacao).Error; err != nil {
		return nil, fmt.Errorf("failed to create transaction: %w", err)
	}

	return transacao, nil
}

func (t *TransacoesRepository) Delete(ctx context.Context, id int) error {
	if err := t.conn.WithContext(ctx).Table(t.table).Where("id = ?", id).Delete(&domain.Transacao{}).Error; err != nil {
		return fmt.Errorf("failed to delete transaction: %w", err)
	}

	return nil
}

func (t *TransacoesRepository) GetAllByUsuarioIDAndMes(ctx context.Context, usuarioID int, mes time.Time) ([]domain.Transacao, error) {
	inicioMes, inicioProximoMes := monthRange(mes)

	var transacoes []domain.Transacao
	if err := t.conn.WithContext(ctx).
		Table(t.table).
		Where("usuario_id = ? AND data >= ? AND data < ?", usuarioID, inicioMes, inicioProximoMes).
		Order("data DESC, id DESC").
		Find(&transacoes).Error; err != nil {
		return nil, fmt.Errorf("failed to get transactions by user id and month: %w", err)
	}

	return transacoes, nil
}

func (t *TransacoesRepository) GetTotalByUsuarioIDAndMesIgnoringCategoriaID(ctx context.Context, usuarioID int, mes time.Time, ignoredCategoriaID int) (float64, error) {
	inicioMes, inicioProximoMes := monthRange(mes)

	var total float64
	if err := t.conn.WithContext(ctx).
		Table(t.table).
		Select("COALESCE(SUM(valor), 0)").
		Where(
			"usuario_id = ? AND data >= ? AND data < ? AND categoria_id <> ?",
			usuarioID,
			inicioMes,
			inicioProximoMes,
			ignoredCategoriaID,
		).
		Scan(&total).Error; err != nil {
		return 0, fmt.Errorf("failed to get total transactions by month: %w", err)
	}

	return total, nil
}

func monthRange(mes time.Time) (time.Time, time.Time) {
	inicioMes := time.Date(mes.Year(), mes.Month(), 1, 0, 0, 0, 0, mes.Location())
	return inicioMes, inicioMes.AddDate(0, 1, 0)
}
