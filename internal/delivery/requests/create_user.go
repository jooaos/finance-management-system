package requests

type CreateUserRequest struct {
	Nome  string `json:"nome" validate:"required"`
	Email string `json:"email" validate:"required,email"`
}
