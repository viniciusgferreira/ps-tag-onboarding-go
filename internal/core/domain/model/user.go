package model

import (
	"errors"
	"regexp"
)

type User struct {
	ID        string `bson:"_id,omitempty" json:"id,omitempty"`
	FirstName string `bson:"firstName" json:"firstName"`
	LastName  string `bson:"lastName" json:"lastName"`
	Email     string `bson:"email" json:"email"`
	Age       int    `bson:"age" json:"age"`
}

func NewUser(id string, firstName string, lastName string, email string, age int) (*User, error) {
	user := &User{
		ID:        id,
		FirstName: firstName,
		LastName:  lastName,
		Email:     email,
		Age:       age,
	}
	err := validateUser(user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func validateUser(u *User) error {
	if len(u.FirstName) == 0 {
		return errors.New("first name is required")
	}

	if len(u.LastName) == 0 {
		return errors.New("last name is required")
	}
	if len(u.Email) == 0 {
		return errors.New("email is required")
	}
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(u.Email) {
		return errors.New("invalid email")
	}
	if u.Age < 18 {
		return errors.New("user must be at least 18 years old")
	}
	return nil
}
