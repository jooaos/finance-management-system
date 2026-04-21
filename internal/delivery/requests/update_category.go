package requests

type UpdateCategoryRequest struct {
	ID        int    `json:"id" validate:"required,gt=0"`
	Nome      string `json:"nome" validate:"required"`
	UsuarioID int    `json:"usuario_id" validate:"required,gt=0"`
}
