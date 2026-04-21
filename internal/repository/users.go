package repository

import (
	"context"
	"fmt"

	"finance-management/internal/domain"

	"gorm.io/gorm"
)

var (
	usersTable = "usuarios"
)

type UserRepository struct {
	conn  *gorm.DB
	table string
}

func NewUserRepository(
	conn *gorm.DB,
) *UserRepository {
	return &UserRepository{
		conn:  conn,
		table: usersTable,
	}
}

func (u *UserRepository) Create(ctx context.Context, usuario *domain.Usuario) (*domain.Usuario, error) {
	if err := u.conn.WithContext(ctx).Table(u.table).Create(usuario).Error; err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return usuario, nil
}

func (u *UserRepository) GetByEmail(ctx context.Context, email string) (*domain.Usuario, error) {
	var usuario domain.Usuario
	if err := u.conn.WithContext(ctx).Table(u.table).Where("email = ?", email).First(&usuario).Error; err != nil {
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}

	return &usuario, nil
}

func (u *UserRepository) GetByID(ctx context.Context, id int) (*domain.Usuario, error) {
	var usuario domain.Usuario
	if err := u.conn.WithContext(ctx).Table(u.table).Where("id = ?", id).First(&usuario).Error; err != nil {
		return nil, fmt.Errorf("failed to get user by id: %w", err)
	}

	return &usuario, nil
}
