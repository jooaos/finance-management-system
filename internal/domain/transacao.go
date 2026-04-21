package domain

import "time"

type Transacao struct {
	ID          int       `json:"id"`
	UsuarioID   int       `json:"usuario_id"`
	CategoriaID int       `json:"categoria_id"`
	Valor       float64   `json:"valor"`
	Data        time.Time `json:"data"`
	Descricao   string    `json:"descricao"`
	Tipo        string    `json:"tipo"`
	Parcelas    int       `json:"parcelas"`
}
