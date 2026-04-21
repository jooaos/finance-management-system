package requests

import "time"

type CreateTransactionRequest struct {
	UsuarioID   int       `json:"usuario_id" validate:"required,gt=0"`
	CategoriaID int       `json:"categoria_id" validate:"required,gt=0"`
	Valor       float64   `json:"valor" validate:"required,gt=0"`
	Data        time.Time `json:"data" validate:"required"`
	Descricao   string    `json:"descricao"`
	Tipo        string    `json:"tipo" validate:"required"`
	Parcelas    int       `json:"parcelas" validate:"gte=0"`
}
