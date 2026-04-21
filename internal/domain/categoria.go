package domain

type Categoria struct {
	ID        int    `json:"id"`
	Nome      string `json:"nome"`
	UsuarioID int    `json:"usuario_id"`
}
