package dto

type InsertAccountRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required,min=4"`
	FullName string `json:"full_name"`
	Gender   string `json:"gender" binding:"max=1"`
}
