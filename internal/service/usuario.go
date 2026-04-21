package service

import (
	"context"
	"errors"
	"fmt"

	"finance-management/internal/delivery/requests"
	"finance-management/internal/domain"
	"finance-management/internal/repository"

	"gorm.io/gorm"
)

var defaultCategoryNames = []string{
	"Alimentação",
	"Transporte",
	"Lazer",
	"Moradia",
	"Receita",
}

type UsuarioService struct {
	userRepository      *repository.UserRepository
	categoriaRepository *repository.CategoriaRepository
}

func NewUsuarioService(
	userRepository *repository.UserRepository,
	categoriaRepository *repository.CategoriaRepository,
) *UsuarioService {
	return &UsuarioService{
		userRepository:      userRepository,
		categoriaRepository: categoriaRepository,
	}
}

func (u *UsuarioService) Create(ctx context.Context, request requests.CreateUserRequest) (*domain.Usuario, error) {
	existingUsuario, err := u.userRepository.GetByEmail(ctx, request.Email)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("failed to check user email: %w", err)
	}

	if existingUsuario != nil {
		return nil, fmt.Errorf("user with email %q already exists", request.Email)
	}

	usuario := &domain.Usuario{
		Nome:  request.Nome,
		Email: request.Email,
	}

	createdUsuario, err := u.userRepository.Create(ctx, usuario)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	for _, categoryName := range defaultCategoryNames {
		categoria := &domain.Categoria{
			Nome:      categoryName,
			UsuarioID: createdUsuario.ID,
		}

		if _, err := u.categoriaRepository.Create(ctx, categoria); err != nil {
			return nil, fmt.Errorf("failed to create default category %q: %w", categoryName, err)
		}
	}

	return createdUsuario, nil
}
