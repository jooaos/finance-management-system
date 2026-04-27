package service

import (
	"context"
	"fmt"
	"math"
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

type CategoryMonthlySpendingSummary struct {
	CategoriaID         int     `json:"categoria_id"`
	CategoriaNome       string  `json:"categoria_nome"`
	Orcamento           float64 `json:"orcamento"`
	Gasto               float64 `json:"gasto"`
	Disponivel          float64 `json:"disponivel"`
	PercentualUtilizado float64 `json:"percentual_utilizado"`
	Excedido            bool    `json:"excedido"`
}

type MonthlyCommitmentProjection struct {
	Mes                    time.Time `json:"mes"`
	TotalReceita           float64   `json:"total_receita"`
	TotalDespesa           float64   `json:"total_despesa"`
	SaldoProjetado         float64   `json:"saldo_projetado"`
	PercentualComprometido float64   `json:"percentual_comprometido"`
}

type TransacaoService struct {
	transacoesRepository *repository.TransacoesRepository
	categoriaRepository  *repository.CategoriaRepository
	orcamentosRepository *repository.OrcamentosRepository
}

func NewTransacaoService(
	transacoesRepository *repository.TransacoesRepository,
	categoriaRepository *repository.CategoriaRepository,
	orcamentosRepository *repository.OrcamentosRepository,
) *TransacaoService {
	return &TransacaoService{
		transacoesRepository: transacoesRepository,
		categoriaRepository:  categoriaRepository,
		orcamentosRepository: orcamentosRepository,
	}
}

func (t *TransacaoService) Create(ctx context.Context, request requests.CreateTransactionRequest) ([]domain.Transacao, error) {
	totalParcelas := request.Parcelas
	if totalParcelas < 1 {
		totalParcelas = 1
	}

	isReceita, err := t.isReceitaCategory(ctx, request.UsuarioID, request.CategoriaID)
	if err != nil {
		return nil, err
	}
	if isReceita {
		totalParcelas = 1
	}

	valoresParcelas := splitInstallments(request.Valor, totalParcelas)

	createdTransacoes := make([]domain.Transacao, 0, totalParcelas)
	for parcela := 0; parcela < totalParcelas; parcela++ {
		transacao := &domain.Transacao{
			UsuarioID:   request.UsuarioID,
			CategoriaID: request.CategoriaID,
			Valor:       valoresParcelas[parcela],
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

func (t *TransacaoService) isReceitaCategory(ctx context.Context, usuarioID int, categoriaID int) (bool, error) {
	categorias, err := t.categoriaRepository.GetAllByUsuarioID(ctx, usuarioID)
	if err != nil {
		return false, fmt.Errorf("failed to get user categories: %w", err)
	}

	for _, categoria := range categorias {
		if categoria.ID == categoriaID {
			return strings.EqualFold(categoria.Nome, receitaCategoryName), nil
		}
	}

	return false, nil
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

func (t *TransacaoService) GetCategoryMonthlySummary(ctx context.Context, usuarioID int, mes time.Time) ([]CategoryMonthlySpendingSummary, error) {
	categorias, err := t.categoriaRepository.GetAllByUsuarioID(ctx, usuarioID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user categories: %w", err)
	}

	orcamentos, err := t.orcamentosRepository.GetByUsuarioIDAndMes(ctx, usuarioID, normalizeMes(mes))
	if err != nil {
		return nil, fmt.Errorf("failed to get monthly budgets: %w", err)
	}

	transacoes, err := t.transacoesRepository.GetAllByUsuarioIDAndMes(ctx, usuarioID, mes)
	if err != nil {
		return nil, fmt.Errorf("failed to get monthly transactions: %w", err)
	}

	categoryTotals := make(map[int]float64, len(categorias))
	categoryBudgets := make(map[int]float64, len(orcamentos))
	categoryNames := make(map[int]string, len(categorias))
	for _, categoria := range categorias {
		categoryNames[categoria.ID] = categoria.Nome
	}
	for _, orcamento := range orcamentos {
		categoryBudgets[orcamento.CategoriaID] = orcamento.Limite
	}

	for _, transacao := range transacoes {
		if strings.EqualFold(categoryNames[transacao.CategoriaID], receitaCategoryName) {
			continue
		}

		categoryTotals[transacao.CategoriaID] += transacao.Valor
	}

	summaries := make([]CategoryMonthlySpendingSummary, 0, len(categorias))
	for _, categoria := range categorias {
		if strings.EqualFold(categoria.Nome, receitaCategoryName) {
			continue
		}

		orcamento := categoryBudgets[categoria.ID]
		gasto := categoryTotals[categoria.ID]
		disponivel := orcamento - gasto
		excedido := gasto > orcamento

		var percentualUtilizado float64
		if orcamento > 0 {
			percentualUtilizado = (gasto / orcamento) * 100
		}

		summaries = append(summaries, CategoryMonthlySpendingSummary{
			CategoriaID:         categoria.ID,
			CategoriaNome:       categoria.Nome,
			Orcamento:           orcamento,
			Gasto:               gasto,
			Disponivel:          disponivel,
			PercentualUtilizado: percentualUtilizado,
			Excedido:            excedido,
		})
	}

	return summaries, nil
}

func (t *TransacaoService) GetCommitmentProjection(ctx context.Context, usuarioID int, mes time.Time, months int) ([]MonthlyCommitmentProjection, error) {
	if months < 1 {
		months = 1
	}

	projections := make([]MonthlyCommitmentProjection, 0, months)
	baseMonth := time.Date(mes.Year(), mes.Month(), 1, 0, 0, 0, 0, mes.Location())

	for monthIndex := 0; monthIndex < months; monthIndex++ {
		currentMonth := baseMonth.AddDate(0, monthIndex, 0)
		totalReceita, totalDespesa, err := t.getMonthlyRevenueAndExpense(ctx, usuarioID, currentMonth)
		if err != nil {
			return nil, err
		}

		var percentualComprometido float64
		if totalReceita > 0 {
			percentualComprometido = (totalDespesa / totalReceita) * 100
		}

		projections = append(projections, MonthlyCommitmentProjection{
			Mes:                    currentMonth,
			TotalReceita:           totalReceita,
			TotalDespesa:           totalDespesa,
			SaldoProjetado:         totalReceita - totalDespesa,
			PercentualComprometido: percentualComprometido,
		})
	}

	return projections, nil
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

func splitInstallments(total float64, parcelas int) []float64 {
	if parcelas <= 1 {
		return []float64{roundToCents(total)}
	}

	totalCentavos := int(math.Round(total * 100))
	baseCentavos := totalCentavos / parcelas
	restanteCentavos := totalCentavos - (baseCentavos * parcelas)

	valores := make([]float64, parcelas)
	for i := 0; i < parcelas; i++ {
		centavos := baseCentavos
		if i == parcelas-1 {
			centavos += restanteCentavos
		}
		valores[i] = float64(centavos) / 100
	}

	return valores
}

func roundToCents(value float64) float64 {
	return math.Round(value*100) / 100
}
