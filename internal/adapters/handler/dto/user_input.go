package dto

type UserInput struct {
	FirstName string `json:"firstName" binding:"required"`
	LastName  string `json:"lastName" binding:"required"`
	Email     string `json:"email" binding:"required,email"`
	Age       int    `json:"age" binding:"required,gte=18" gte:"user must be at least 18 years old"`
}
