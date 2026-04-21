package service

import (
	"context"
	"fmt"
	"strings"
	"time"

	"finance-management/internal/delivery/requests"
	"finance-management/internal/domain"
	"finance-management/internal/repository"
)

const receitaCategoryName = "Receita"

type MonthlyTransactionSummary struct {
	TotalReceita float64 `json:"total_receita"`
	TotalDespesa float64 `json:"total_despesa"`
	SaldoAtual   float64 `json:"saldo_atual"`
}

type MonthlySpendingProgress struct {
	TotalGastoMes     float64 `json:"total_gasto_mes"`
	SaldoTotal        float64 `json:"saldo_total"`
	PercentualDoSaldo float64 `json:"percentual_do_saldo"`
}

type TransacaoService struct {
	transacoesRepository *repository.TransacoesRepository
	categoriaRepository  *repository.CategoriaRepository
}

func NewTransacaoService(
	transacoesRepository *repository.TransacoesRepository,
	categoriaRepository *repository.CategoriaRepository,
) *TransacaoService {
	return &TransacaoService{
		transacoesRepository: transacoesRepository,
		categoriaRepository:  categoriaRepository,
	}
}

func (t *TransacaoService) Create(ctx context.Context, request requests.CreateTransactionRequest) ([]domain.Transacao, error) {
	totalParcelas := request.Parcelas
	if totalParcelas < 1 {
		totalParcelas = 1
	}

	createdTransacoes := make([]domain.Transacao, 0, totalParcelas)
	for parcela := 0; parcela < totalParcelas; parcela++ {
		transacao := &domain.Transacao{
			UsuarioID:   request.UsuarioID,
			CategoriaID: request.CategoriaID,
			Valor:       request.Valor,
			Data:        request.Data.AddDate(0, parcela, 0),
			Descricao:   request.Descricao,
			Tipo:        request.Tipo,
			Parcelas:    totalParcelas,
		}

		createdTransacao, err := t.transacoesRepository.Create(ctx, transacao)
		if err != nil {
			return nil, fmt.Errorf("failed to create transaction installment %d/%d: %w", parcela+1, totalParcelas, err)
		}

		createdTransacoes = append(createdTransacoes, *createdTransacao)
	}

	return createdTransacoes, nil
}

func (t *TransacaoService) GetMonthlySummary(ctx context.Context, usuarioID int, mes time.Time) (*MonthlyTransactionSummary, error) {
	totalReceita, totalDespesa, err := t.getMonthlyRevenueAndExpense(ctx, usuarioID, mes)
	if err != nil {
		return nil, err
	}

	return &MonthlyTransactionSummary{
		TotalReceita: totalReceita,
		TotalDespesa: totalDespesa,
		SaldoAtual:   totalReceita - totalDespesa,
	}, nil
}

func (t *TransacaoService) GetMonthlySpendingProgress(ctx context.Context, usuarioID int, mes time.Time) (*MonthlySpendingProgress, error) {
	totalReceita, totalDespesa, err := t.getMonthlyRevenueAndExpense(ctx, usuarioID, mes)
	if err != nil {
		return nil, err
	}

	var percentualDoSaldo float64
	if totalReceita > 0 {
		percentualDoSaldo = (totalDespesa / totalReceita) * 100
	}

	return &MonthlySpendingProgress{
		TotalGastoMes:     totalDespesa,
		SaldoTotal:        totalReceita,
		PercentualDoSaldo: percentualDoSaldo,
	}, nil
}

func (t *TransacaoService) getMonthlyRevenueAndExpense(ctx context.Context, usuarioID int, mes time.Time) (float64, float64, error) {
	categorias, err := t.categoriaRepository.GetAllByUsuarioID(ctx, usuarioID)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to get user categories: %w", err)
	}

	categoryNames := make(map[int]string, len(categorias))
	for _, categoria := range categorias {
		categoryNames[categoria.ID] = categoria.Nome
	}

	transacoes, err := t.transacoesRepository.GetAllByUsuarioIDAndMes(ctx, usuarioID, mes)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to get monthly transactions: %w", err)
	}

	var totalReceita float64
	var totalDespesa float64
	for _, transacao := range transacoes {
		if strings.EqualFold(categoryNames[transacao.CategoriaID], receitaCategoryName) {
			totalReceita += transacao.Valor
			continue
		}

		totalDespesa += transacao.Valor
	}

	return totalReceita, totalDespesa, nil
}
