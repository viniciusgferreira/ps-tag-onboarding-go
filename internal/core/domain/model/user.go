package model

import "github.com/viniciusgferreira/ps-tag-onboarding-go/internal/adapters/handler/dto"

type User struct {
	ID        string `bson:"_id,omitempty" json:"id,omitempty"`
	FirstName string `bson:"firstName" json:"firstName"`
	LastName  string `bson:"lastName" json:"lastName"`
	Email     string `bson:"email" json:"email"`
	Age       int    `bson:"age" json:"age"`
}

func NewUser(dto dto.UserInput) User {
	return User{
		FirstName: dto.FirstName,
		LastName:  dto.LastName,
		Email:     dto.Email,
		Age:       dto.Age,
	}
}
