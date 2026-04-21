package domain

import "time"

type Orcamento struct {
	ID          int       `json:"id"`
	UsuarioID   int       `json:"usuario_id"`
	CategoriaID int       `json:"categoria_id"`
	Limite      float64   `json:"limite"`
	Mes         time.Time `json:"mes"`
}
