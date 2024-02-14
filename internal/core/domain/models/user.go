package models

import (
	"errors"
	"github.com/google/uuid"
	"regexp"
)

type User struct {
	id        string
	firstName string
	lastName  string
	email     string
	age       int
}

func NewUser(id, firstName, lastName, email string, age int) User {
	user := User{
		id:        id,
		firstName: firstName,
		lastName:  lastName,
		email:     email,
		age:       age,
	}
	return user
}
func (u *User) validate() error {
	if err := u.validateID(); err != nil {
		return err
	}

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
	if len(u.id) <= 0 {
		return errors.New("id cannot be empty")
	}
	if err := uuid.Validate(u.id); err != nil {
		return errors.New("invalid id")
	}
	return nil
}

func (u *User) validateName() error {
	if len(u.firstName) <= 0 || len(u.lastName) <= 0 {
		return errors.New("first and last name cannot be empty")
	}
	return nil
}

func (u *User) validateEmail() error {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(u.email) {
		return errors.New("invalid email format")
	}
	return nil
}

func (u *User) validateAge() error {
	if u.age < 18 {
		return errors.New("user must be at least 18 years old")
	}
	return nil
}
