package repository

import (
	"context"
	"fmt"

	"finance-management/internal/domain"

	"gorm.io/gorm"
)

var (
	categoriaTable = "categorias"
)

type CategoriaRepository struct {
	conn  *gorm.DB
	table string
}

func NewCategoriaRepository(
	conn *gorm.DB,
) *CategoriaRepository {
	return &CategoriaRepository{
		conn:  conn,
		table: categoriaTable,
	}
}

func (c *CategoriaRepository) Create(ctx context.Context, categoria *domain.Categoria) (*domain.Categoria, error) {
	if err := c.conn.WithContext(ctx).Table(c.table).Create(categoria).Error; err != nil {
		return nil, fmt.Errorf("failed to create category: %w", err)
	}

	return categoria, nil
}

func (c *CategoriaRepository) GetAllByUsuarioID(ctx context.Context, usuarioID int) ([]domain.Categoria, error) {
	var categorias []domain.Categoria
	if err := c.conn.WithContext(ctx).Table(c.table).Where("usuario_id = ?", usuarioID).Find(&categorias).Error; err != nil {
		return nil, fmt.Errorf("failed to get categories by user id: %w", err)
	}

	return categorias, nil
}

func (c *CategoriaRepository) Delete(ctx context.Context, id int) error {
	if err := c.conn.WithContext(ctx).Table(c.table).Where("id = ?", id).Delete(&domain.Categoria{}).Error; err != nil {
		return fmt.Errorf("failed to delete category: %w", err)
	}

	return nil
}

func (c *CategoriaRepository) Update(ctx context.Context, categoria *domain.Categoria) (*domain.Categoria, error) {
	if err := c.conn.WithContext(ctx).Table(c.table).Where("id = ?", categoria.ID).Updates(categoria).Error; err != nil {
		return nil, fmt.Errorf("failed to update category: %w", err)
	}

	return c.GetByID(ctx, categoria.ID)
}

func (c *CategoriaRepository) GetByID(ctx context.Context, id int) (*domain.Categoria, error) {
	var categoria domain.Categoria
	if err := c.conn.WithContext(ctx).Table(c.table).Where("id = ?", id).First(&categoria).Error; err != nil {
		return nil, fmt.Errorf("failed to get category by id: %w", err)
	}

	return &categoria, nil
}
