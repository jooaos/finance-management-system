package requests

type CreateCategoryRequest struct {
	Nome      string `json:"nome" validate:"required"`
	UsuarioID int    `json:"usuario_id" validate:"required,gt=0"`
}
