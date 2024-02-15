package models

import (
	"errors"
	"github.com/google/uuid"
	"regexp"
)

type User struct {
	ID        string `json:"id,omitempty"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Age       int    `json:"age"`
}

func NewUser(id, firstName, lastName, email string, age int) User {
	user := User{
		ID:        id,
		FirstName: firstName,
		LastName:  lastName,
		Email:     email,
		Age:       age,
	}
	return user
}
func (u *User) Validate() error {
	if err := u.validateName(); err != nil {
		return err
	}
	if err := u.validateEmail(); err != nil {
		return err
	}
	if err := u.validateAge(); err != nil {
		return err
	}
	return nil
}

func (u *User) validateID() error {
	if len(u.ID) <= 0 {
		return errors.New("id cannot be empty")
	}
	if err := uuid.Validate(u.ID); err != nil {
		return errors.New("invalid id")
	}
	return nil
}

func (u *User) validateName() error {
	if len(u.FirstName) <= 0 || len(u.LastName) <= 0 {
		return errors.New("first and last name cannot be empty")
	}
	return nil
}

func (u *User) validateEmail() error {
	if len(u.Email) <= 0 {
		return errors.New("email cannot be empty")
	}
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(u.Email) {
		return errors.New("invalid email format")
	}
	return nil
}

func (u *User) validateAge() error {
	if u.Age < 18 {
		return errors.New("user must be at least 18 years old")
	}
	return nil
}
