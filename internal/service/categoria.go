package service

import (
	"context"
	"fmt"

	"finance-management/internal/delivery/requests"
	"finance-management/internal/domain"
	"finance-management/internal/repository"
)

type CategoriaService struct {
	categoriaRepository *repository.CategoriaRepository
}

func NewCategoriaService(categoriaRepository *repository.CategoriaRepository) *CategoriaService {
	return &CategoriaService{
		categoriaRepository: categoriaRepository,
	}
}

func (c *CategoriaService) Create(ctx context.Context, request requests.CreateCategoryRequest) (*domain.Categoria, error) {
	categoria := &domain.Categoria{
		Nome:      request.Nome,
		UsuarioID: request.UsuarioID,
	}

	createdCategoria, err := c.categoriaRepository.Create(ctx, categoria)
	if err != nil {
		return nil, fmt.Errorf("failed to create category: %w", err)
	}

	return createdCategoria, nil
}

func (c *CategoriaService) GetAllByUsuarioID(ctx context.Context, usuarioID int) ([]domain.Categoria, error) {
	categorias, err := c.categoriaRepository.GetAllByUsuarioID(ctx, usuarioID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user categories: %w", err)
	}

	return categorias, nil
}

func (c *CategoriaService) Update(ctx context.Context, request requests.UpdateCategoryRequest) (*domain.Categoria, error) {
	categoria := &domain.Categoria{
		ID:        request.ID,
		Nome:      request.Nome,
		UsuarioID: request.UsuarioID,
	}

	updatedCategoria, err := c.categoriaRepository.Update(ctx, categoria)
	if err != nil {
		return nil, fmt.Errorf("failed to update category: %w", err)
	}

	return updatedCategoria, nil
}

func (c *CategoriaService) Delete(ctx context.Context, id int) error {
	if err := c.categoriaRepository.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete category: %w", err)
	}

	return nil
}

func (c *CategoriaService) GetByID(ctx context.Context, id int) (*domain.Categoria, error) {
	categoria, err := c.categoriaRepository.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get category by id: %w", err)
	}

	return categoria, nil
}
