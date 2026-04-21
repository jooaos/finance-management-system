package requests

import "time"

type UpdateBudgetRequest struct {
	ID          int       `json:"id" validate:"required,gt=0"`
	UsuarioID   int       `json:"usuario_id" validate:"required,gt=0"`
	CategoriaID int       `json:"categoria_id" validate:"required,gt=0"`
	Limite      float64   `json:"limite" validate:"required,gt=0"`
	Mes         time.Time `json:"mes" validate:"required"`
}
